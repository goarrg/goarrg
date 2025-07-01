/*
Copyright 2025 The goARRG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package golang

import (
	"fmt"
	"go/types"
	"io"
	"os"
	"slices"
	"strings"
	"unsafe"

	"goarrg.com/debug"
	"goarrg.com/toolchain"

	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type exportedStructField struct {
	Name, CType, GoType string
	ArraySize           int64
}

type exportedType struct {
	Def             string
	GoType          string
	Fields          []exportedStructField
	ArrayElemCType  string
	ArrayElemGoType string
	ArraySize       int64
	Depends         []string
	IsHandle        bool
	Visited         bool
}

type exportedFuncArg struct {
	Name   string
	CType  string
	GoType string
	GoElem string
}

type exportedFunc struct {
	TargetName string
	Recv       string
	Args       []exportedFuncArg
	IsVariadic bool
	Return     exportedFuncArg
}

func basicTypeToCType(t types.BasicKind) string {
	switch t {
	case types.Bool:
		return "uint32_t"
	case types.Int:
		switch unsafe.Sizeof(int(0)) {
		case unsafe.Sizeof(int64(0)):
			return "int64_t"
		case unsafe.Sizeof(int32(0)):
			return "int32_t"
		}
	case types.Int8:
		return "int8_t"
	case types.Int16:
		return "int16_t"
	case types.Int32:
		return "int32_t"
	case types.Int64:
		return "int64_t"
	case types.Uint:
		switch unsafe.Sizeof(uint(0)) {
		case unsafe.Sizeof(uint64(0)):
			return "uint64_t"
		case unsafe.Sizeof(uint32(0)):
			return "uint32_t"
		}
	case types.Uint8:
		return "uint8_t"
	case types.Uint16:
		return "uint16_t"
	case types.Uint32:
		return "uint32_t"
	case types.Uint64:
		return "uint64_t"
	case types.Uintptr:
		return "uintptr_t"
	case types.Float32:
		return "float"
	case types.Float64:
		return "double"
	case types.String:
		return "char*"
	}
	panic(fmt.Sprintf("Unhandled type: %d", t))
}

func goNameToC(pkg, name string) string {
	return "Go_" + pkg + "_" + name
}

func exportGoFunc(mF map[string]*exportedFunc, mT map[string]*exportedType, mI map[string]struct{}, o *types.Func) {
	dName := goNameToC(o.Pkg().Name(), o.Name())
	if mF[dName] != nil {
		return
	}
	sig := o.Type().(*types.Signature)
	ex := exportedFunc{
		TargetName: o.Pkg().Name() + "." + o.Name(),
		IsVariadic: sig.Variadic(),
	}
	if sig.Recv() != nil {
		name, _ := exportGoTypeName(mI, sig.Recv().Type())
		ex.Recv = goTypeName(sig.Recv().Type())
		ex.TargetName = o.Name()
		dName = name + "_" + o.Name()
		ex.Args = append(ex.Args, exportedFuncArg{
			Name:   "recv",
			CType:  name,
			GoType: ex.Recv,
		})
	}
	p := sig.Params()
	complete := true
	for i := range p.Len() {
		pName, isBasic := exportGoTypeName(mI, p.At(i).Type())
		if pName == "" {
			complete = false
			break
		}
		if !isBasic {
			exportGoType(mT, mI, pName, p.At(i).Type())
			{
				_, isPointer := p.At(i).Type().Underlying().(*types.Pointer)
				_, isInterface := p.At(i).Type().Underlying().(*types.Interface)
				if mT[pName].IsHandle && !(isPointer || isInterface) {
					complete = false
					break
				}
			}
		}
		ex.Args = append(ex.Args, exportedFuncArg{
			Name:   p.At(i).Name(),
			CType:  pName,
			GoType: goTypeName(p.At(i).Type()),
		})
	}
	if !complete {
		debug.WPrintf("UNHANDLED: unhandled args: %s", ex.TargetName)
		return
	}
	if sig.Results().Len() > 1 {
		debug.WPrintf("UNHANDLED: funcs with more than one return: %s", ex.TargetName)
		return
	}
	if sig.Results().Len() == 1 {
		ret := sig.Results().At(0).Type()
		pName, isBasic := exportGoTypeName(mI, ret)
		if pName == "" {
			debug.WPrintf("UNHANDLED: return type: %s", ex.TargetName)
			return
		}
		if strings.HasSuffix(pName, "_Slice") {
			debug.WPrintf("UNHANDLED: slice return type: %s", ex.TargetName)
			return
		}
		if !isBasic {
			exportGoType(mT, mI, pName, sig.Results().At(0).Type())
			for _, f := range mT[pName].Fields {
				if strings.HasSuffix(f.CType, "_Slice") {
					debug.WPrintf("UNHANDLED: slice return type: %s", ex.TargetName)
					return
				}
			}
		} else if pName == "char*" {
			debug.WPrintf("UNHANDLED: string return type: %s", ex.TargetName)
			return
		}
		ex.Return = exportedFuncArg{
			CType:  pName,
			GoType: goTypeName(ret),
		}
	}
	mF[dName] = &ex
}

func exportGoStruct(mT map[string]*exportedType, mI map[string]struct{}, name string, o types.Type) {
	if mT[name] != nil {
		return
	}
	ex := exportedType{
		Def:    "typedef struct {\n",
		GoType: goTypeName(o),
	}
	t := o.Underlying().(*types.Struct)
	complete := t.NumFields() > 0
	for i := range t.NumFields() {
		f := t.Field(i)
		if strings.HasPrefix(f.Name(), "_") {
			continue
		}
		if !f.Exported() {
			complete = false
			break
		}
		dName, _ := exportGoTypeName(mI, f.Type())
		if dName == "" {
			complete = false
			break
		}
	}
	if complete {
		for i := range t.NumFields() {
			f := t.Field(i)
			if strings.HasPrefix(f.Name(), "_") {
				continue
			}
			dName, isBasic := exportGoTypeName(mI, f.Type())
			if arr, ok := f.Type().(*types.Array); ok {
				if !isBasic {
					ex.Depends = append(ex.Depends, dName)
					exportGoType(mT, mI, dName, arr.Elem())
				}
				ex.Def += fmt.Sprintf("\t%s %s[%d];\n", dName, f.Name(), arr.Len())
				ex.Fields = append(ex.Fields, exportedStructField{
					Name:      f.Name(),
					CType:     dName,
					GoType:    goTypeName(f.Type()),
					ArraySize: arr.Len(),
				})
			} else {
				ex.Depends = append(ex.Depends, dName)
				if !isBasic {
					exportGoType(mT, mI, dName, f.Type())
				}
				ex.Def += fmt.Sprintf("\t%s %s;\n", dName, f.Name())
				ex.Fields = append(ex.Fields, exportedStructField{
					Name:   f.Name(),
					CType:  dName,
					GoType: goTypeName(f.Type()),
				})
			}
		}
	}
	slices.Sort(ex.Depends)
	ex.Depends = slices.Compact(ex.Depends)
	if !complete {
		ex.Def = "GO_HANDLE(" + name + ");"
		ex.GoType = "*" + ex.GoType
		ex.IsHandle = true
	}
	mT[name] = &ex
}

func exportGoType(m map[string]*exportedType, mI map[string]struct{}, name string, o types.Type) {
	if m[name] != nil {
		return
	}
	gType := goTypeName(o)
	switch underlying := o.Underlying().(type) {
	case *types.Basic:
		m[name] = &exportedType{
			Def:    "typedef " + basicTypeToCType(underlying.Kind()),
			GoType: gType,
		}
	case *types.Struct:
		exportGoStruct(m, mI, name, o)
	case *types.Interface:
		m[name] = &exportedType{
			Def:      "GO_HANDLE(" + name + ");",
			GoType:   gType,
			IsHandle: true,
		}
	case *types.Pointer:
		m[name] = &exportedType{
			Def:      "GO_HANDLE(" + name + ");",
			GoType:   gType,
			IsHandle: true,
		}
	case *types.Slice:
		tName, isBasic := exportGoTypeName(mI, underlying.Elem())
		m[name] = &exportedType{
			Def: "typedef struct {\n" +
				"\t" + tName + "* ptr;\n" +
				"\tsize_t len;\n" +
				"}",
			GoType: gType,
		}
		if !isBasic {
			exportGoType(m, mI, tName, underlying.Elem())
			m[name].Depends = append(m[name].Depends, tName)
		}
	case *types.Array:
		tName, isBasic := exportGoTypeName(mI, underlying.Elem())
		m[name] = &exportedType{
			Def:             fmt.Sprintf("typedef %s %s[%d];", tName, name, underlying.Len()),
			GoType:          gType,
			ArraySize:       underlying.Len(),
			ArrayElemCType:  tName,
			ArrayElemGoType: goTypeName(underlying.Elem()),
		}
		if !isBasic {
			exportGoType(m, mI, tName, underlying.Elem())
			m[name].Depends = append(m[name].Depends, tName)
		}
	case *types.Map:
		/*
			case *types.Map:
					kName, _ := exportGoTypeName(mI, underlying.Key())
					tName, _ := exportGoTypeName(mI, underlying.Elem())
					exportGoUnderlyingType(m, mI, kName, underlying.Key().Underlying())
					exportGoUnderlyingType(m, mI, tName, underlying.Elem().Underlying())
					m[name] = &exportedType{
						Def:      "GO_HANDLE(" + name + ");",
						IsHandle: true,
					}
		*/
	default:
		panic(fmt.Sprintf("UNHANDLED: %s %T %T\n", name, o, underlying))
	}
}

func goTypeName(t types.Type) string {
	switch kind := t.(type) {
	case *types.Alias:
		return kind.Obj().Pkg().Name() + "." + kind.Obj().Name()
	case *types.Named:
		if kind.Origin().Obj().Pkg() == nil {
			return ""
		}
		return kind.Obj().Pkg().Name() + "." + kind.Obj().Name()
	case *types.Basic:
		return kind.Name()
	case *types.Pointer:
		name := goTypeName(kind.Elem())
		return "*" + name
	case *types.Slice:
		name := goTypeName(kind.Elem())
		return "[]" + name
	case *types.Array:
		name := goTypeName(kind.Elem())
		return fmt.Sprintf("[%d]%s", kind.Len(), name)
	default:
		panic(fmt.Sprintf("UNHANDLED: %T\n", t))
	}
}

func exportGoTypeName(mI map[string]struct{}, t types.Type) (string, bool) {
	switch kind := t.(type) {
	case *types.Alias:
		if kind.Origin().Obj().Pkg() == nil {
			return "", false
		}
		mI[kind.Origin().Obj().Pkg().Path()] = struct{}{}
		_, isBasic := kind.Underlying().(*types.Basic)
		return goNameToC(kind.Obj().Pkg().Name(), kind.Obj().Name()), isBasic
	case *types.Named:
		if kind.Origin().Obj().Pkg() == nil {
			return "", false
		}
		mI[kind.Origin().Obj().Pkg().Path()] = struct{}{}
		_, isBasic := kind.Underlying().(*types.Basic)
		return goNameToC(kind.Obj().Pkg().Name(), kind.Obj().Name()), isBasic
	case *types.Basic:
		if kind.Kind() == types.Bool {
			return "Go_Boolean32", true
		} else {
			return basicTypeToCType(kind.Kind()), true
		}
	case *types.Pointer:
		name, isBasic := exportGoTypeName(mI, kind.Elem())
		if isBasic {
			name = "Go_" + name
		}
		return name, false
	case *types.Slice:
		name, isBasic := exportGoTypeName(mI, kind.Elem())
		if name == "char*" {
			name = "String"
		}
		if isBasic {
			name = "Go_" + name
		}
		return name + "_Slice", false
	case *types.Map:
		debug.WPrintf("UNHANDLED: map type: %s", kind)
		return "", false
	case *types.Struct:
		debug.WPrintf("UNHANDLED: anon struct type: %s", kind)
		return "", false
	case *types.Array:
		return exportGoTypeName(mI, kind.Elem())
	default:
		panic(fmt.Sprintf("UNHANDLED: %T %s\n", t, t))
	}
}

func writeExportedType(m map[string]*exportedType, out io.Writer, name string) {
	for _, d := range m[name].Depends {
		if m[d] != nil && !m[d].Visited {
			writeExportedType(m, out, d)
		}
	}
	if m[name] != nil && !m[name].Visited {
		if m[name].IsHandle || m[name].ArraySize > 0 {
			fmt.Fprintf(out, "%s\n", m[name].Def)
		} else {
			if strings.Contains(m[name].Def, "{") && !strings.HasSuffix(m[name].Def, "}") {
				m[name].Def += "}"
			}
			fmt.Fprintf(out, "%s %s;\n", m[name].Def, name)
		}
		m[name].Visited = true
	}
}

func writeExportedTypeConversion(m map[string]*exportedType, out io.Writer, ctype string) {
	t := m[ctype]
	convert := func(arg, ctype, gotype string) string {
		if m[ctype] == nil {
			switch gotype {
			case "string":
				return fmt.Sprintf("C.GoString(%s)", arg)
			case "bool":
				return fmt.Sprintf("convert_from_Go_Boolean32(%s)", arg)
			default:
				return fmt.Sprintf("%s(%s)", gotype, arg)
			}
		}
		return fmt.Sprintf("convert_from_%s(%s)", ctype, arg)
	}
	for _, d := range t.Depends {
		if m[d] != nil && !m[d].Visited {
			writeExportedTypeConversion(m, out, d)
		}
	}
	if !t.Visited {
		fmt.Fprintf(out, "func convert_from_%[1]s(arg C.%[1]s) %[2]s {\n", ctype, t.GoType)
		switch {
		case len(t.Fields) > 0:
			fmt.Fprintf(out, "\treturn %s {\n", t.GoType)
			for _, f := range t.Fields {
				if f.ArraySize > 0 {
					fmt.Fprintf(out, "\t\t%s: %s{\n", f.Name, f.GoType)
					for i := range f.ArraySize {
						fmt.Fprintf(out, "\t\t\t%s,\n", convert(fmt.Sprintf("arg.%s[%d]", f.Name, i), f.CType, f.GoType))
					}
					fmt.Fprintf(out, "\t\t},\n")
				} else {
					fmt.Fprintf(out, "\t\t%s: %s,\n",
						f.Name, convert("arg."+f.Name, f.CType, f.GoType))
				}
			}
			fmt.Fprintf(out, "\t}\n")
		case strings.HasSuffix(ctype, "_Slice"):
			fmt.Fprintf(out, "\tcgoSlice := unsafe.Slice(arg.ptr, arg.len)\n")
			fmt.Fprintf(out, "\tgoSlice := make(%s, len(cgoSlice))\n", t.GoType)
			fmt.Fprintf(out, "\tfor i, arg := range cgoSlice {\n")
			fmt.Fprintf(out, "\t\tgoSlice[i] = %s\n", convert("arg", strings.TrimSuffix(ctype, "_Slice"), strings.TrimPrefix(t.GoType, "[]")))
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\treturn goSlice\n")
		case t.ArraySize > 0:
			fmt.Fprintf(out, "\tcgoSlice := arg[:]\n")
			fmt.Fprintf(out, "\tgoArr := %s{}\n", t.GoType)
			fmt.Fprintf(out, "\tfor i, arg := range cgoSlice {\n")
			fmt.Fprintf(out, "\t\tgoArr[i] = %s\n", convert("arg", t.ArrayElemCType, t.ArrayElemGoType))
			fmt.Fprintf(out, "\t}\n")
			fmt.Fprintf(out, "\treturn goArr\n")
		case t.IsHandle:
			fmt.Fprintf(out, "\treturn cgo.Handle(uintptr(unsafe.Pointer(arg))).Value().(%s)\n", t.GoType)
		default:
			fmt.Fprintf(out, "\treturn %s(arg)\n", t.GoType)
		}
		fmt.Fprintf(out, "}\n")
		m[ctype].Visited = true
	}
}

func writeOutputTypeConversion(m map[string]*exportedType, out io.Writer, ctype string) {
	t := m[ctype]
	if strings.HasSuffix(ctype, "_Slice") || t == nil || t.ArraySize > 0 {
		return
	}
	convert := func(arg, ctype, gotype string) string {
		if m[ctype] == nil {
			switch gotype {
			case "string":
				return fmt.Sprintf("C.CString(%s)", arg)
			case "bool":
				return fmt.Sprintf("convert_to_Go_Boolean32(%s)", arg)
			default:
				return fmt.Sprintf("C.%s(%s)", ctype, arg)
			}
		}
		return fmt.Sprintf("convert_to_%s(%s)", ctype, arg)
	}
	for _, d := range t.Depends {
		if m[d] != nil && !m[d].Visited {
			writeOutputTypeConversion(m, out, d)
		}
	}
	if !t.Visited {
		fmt.Fprintf(out, "func convert_to_%[1]s(arg %[2]s) C.%[1]s {\n", ctype, t.GoType)
		switch {
		case len(t.Fields) > 0:
			fmt.Fprintf(out, "\treturn C.%s {\n", ctype)
			for _, f := range t.Fields {
				if f.ArraySize > 0 {
					fmt.Fprintf(out, "\t\t%s: %s{\n", f.Name, f.GoType)
					for i := range f.ArraySize {
						fmt.Fprintf(out, "\t\t\t%s,\n", convert(fmt.Sprintf("arg.%s[%d]", f.Name, i), f.CType, f.GoType))
					}
					fmt.Fprintf(out, "\t\t},\n")
				} else {
					fmt.Fprintf(out, "\t\t%s: %s,\n",
						f.Name, convert("arg."+f.Name, f.CType, f.GoType))
				}
			}
			fmt.Fprintf(out, "\t}\n")
		case t.IsHandle:
			fmt.Fprintf(out, "\treturn (C.%s)(unsafe.Pointer(cgo.NewHandle(arg)))\n", ctype)
		default:
			fmt.Fprintf(out, "\treturn C.%s(arg)\n", ctype)
		}
		fmt.Fprintf(out, "}\n")
		m[ctype].Visited = true
	}
}

/*
GenerateCExportFile will parse the packages given and generate C bindings,
it has multiple limitations and doesn't support all types so use with caution.
A quirk of the handle system is that types defined by GO_HANDLE(...)
has to call at most 2 destructors, one to destroy the object itself if needed
and Go_DestroyHandle to destroy the handle.
*/
func GenerateCExportFile(preamble, outfile string, buildflags []string, packagePath ...string) {
	headerOut, err := os.Create(outfile + ".go")
	if err != nil {
		panic(debug.ErrorWrapf(err, "Failed to create file: %s", err))
	}
	{
		defer func() {
			out, err := imports.Process(outfile+".go", nil, nil)
			if err != nil {
				panic(err)
			}
			os.WriteFile(outfile+".go", out, 0o655)
		}()
		defer headerOut.Close()
		defer headerOut.Sync()
	}

	headerOut.WriteString(strings.TrimSpace(preamble) + "\n\n")
	headerOut.WriteString("package main\n\n")
	headerOut.WriteString("/*\n")
	headerOut.WriteString("#include <stddef.h>\n")
	headerOut.WriteString("#include <stdint.h>\n")
	headerOut.WriteString("\n")
	headerOut.WriteString("#define GO_HANDLE(object) typedef struct object##_t* object\n\n")
	headerOut.WriteString("typedef uint32_t Go_Boolean32;\n\n")

	typeMap := map[string]*exportedType{}
	funcMap := map[string]*exportedFunc{}
	importMap := map[string]struct{}{}

	for _, pattern := range packagePath {
		list, err := packages.Load(&packages.Config{
			Mode: packages.NeedName | packages.NeedTypes, BuildFlags: buildflags, Dir: toolchain.WorkingModuleDir(),
		}, pattern)
		if err != nil {
			panic(debug.ErrorWrapf(err, "Failed to load package pattern: %s", packagePath))
		}
		if len(list) == 0 {
			panic(debug.Errorf("No go package returned for: %s", packagePath))
		}

		for _, pkg := range list {
			t := pkg.Types.Scope()
			for _, n := range t.Names() {
				o := t.Lookup(n)
				if o.Exported() {
					switch kind := o.(type) {
					case *types.TypeName:
						dName, _ := exportGoTypeName(importMap, o.Type())
						exportGoType(typeMap, importMap, dName, o.Type())

						set := types.NewMethodSet(o.Type())
						for i := range set.Len() {
							m := set.At(i)
							if m.Obj().Exported() {
								exportGoFunc(funcMap, typeMap, importMap, m.Obj().(*types.Func))
							}
						}
					case *types.Const:
						dName, _ := exportGoTypeName(importMap, o.Type())
						if typeMap[dName] != nil {
							if !strings.HasPrefix(typeMap[dName].Def, "typedef enum") {
								typeMap[dName].Def = "typedef enum {\n"
							}
							typeMap[dName].Def += fmt.Sprintf("\t%s = %s,\n", goNameToC(o.Pkg().Name(), n), kind.Val().ExactString())
						}
						// fmt.Printf("%s %s %v %s\n", n, kind.Type().(*types.Named).Obj().Name(), kind.Val(), o.Type().Underlying().(*types.Basic).Name())
					case *types.Func:
						exportGoFunc(funcMap, typeMap, importMap, kind)
					default:
						panic(fmt.Sprintf("UNHANDLED: %s %T %T\n", n, o, o.Type().Underlying()))
					}
				}
			}
		}
	}

	{
		keys := maps.Keys(typeMap)
		slices.Sort(keys)

		for _, k := range keys {
			writeExportedType(typeMap, headerOut, k)
		}
		headerOut.WriteString("*/\n")
		headerOut.WriteString("import \"C\"\n")
	}
	{
		headerOut.WriteString("import (\n")
		headerOut.WriteString("\"unsafe\"\n")
		headerOut.WriteString("\"runtime/cgo\"\n")
		keys := maps.Keys(importMap)
		slices.Sort(keys)
		for _, k := range keys {
			fmt.Fprintf(headerOut, "%q\n", k)
		}
		headerOut.WriteString(")\n\n")
	}
	{
		keys := maps.Keys(typeMap)
		slices.Sort(keys)

		for _, k := range keys {
			typeMap[k].Visited = false
		}
		for _, k := range keys {
			writeExportedTypeConversion(typeMap, headerOut, k)
		}
		for _, k := range keys {
			typeMap[k].Visited = false
		}

		{
			fmt.Fprintf(headerOut, "func convert_from_Go_Boolean32(value C.Go_Boolean32) bool{\n")
			fmt.Fprintf(headerOut, "\tif value > 0 {\n")
			fmt.Fprintf(headerOut, "\t\treturn true\n")
			fmt.Fprintf(headerOut, "\t} else {\n")
			fmt.Fprintf(headerOut, "\t\treturn false\n")
			fmt.Fprintf(headerOut, "\t}\n")
			fmt.Fprintf(headerOut, "}\n")

			fmt.Fprintf(headerOut, "func convert_to_Go_Boolean32(value bool) C.Go_Boolean32 {\n")
			fmt.Fprintf(headerOut, "\tif value {\n")
			fmt.Fprintf(headerOut, "\t\treturn 1\n")
			fmt.Fprintf(headerOut, "\t} else {\n")
			fmt.Fprintf(headerOut, "\t\treturn 0\n")
			fmt.Fprintf(headerOut, "\t}\n")
			fmt.Fprintf(headerOut, "}\n")
		}
	}
	{
		fmt.Fprintf(headerOut, "// export Go_DestroyHandle\n")
		fmt.Fprintf(headerOut, "func Go_DestroyHandle(handle unsafe.Pointer) {\n")
		fmt.Fprintf(headerOut, "\tcgo.Handle(uintptr(handle)).Delete()\n")
		fmt.Fprintf(headerOut, "}\n")
	}
	{
		keys := maps.Keys(funcMap)
		slices.Sort(keys)
		convertArgs := func(f *exportedFunc) string {
			args := ""
			for _, arg := range f.Args {
				args += "go" + arg.Name + ","
				t := typeMap[arg.CType]
				if t == nil {
					switch arg.GoType {
					case "string":
						fmt.Fprintf(headerOut, "\tgo%[1]s := C.GoString(c%[1]s)\n", arg.Name)
					case "bool":
						fmt.Fprintf(headerOut, "\tgo%[1]s := convert_from_Go_Boolean32(c%[1]s)\n", arg.Name)
					default:
						fmt.Fprintf(headerOut, "\tgo%[1]s := %[2]s(c%[1]s)\n", arg.Name, arg.GoType)
					}
					continue
				}
				fmt.Fprintf(headerOut, "\tgo%[1]s := convert_from_%[2]s(c%[1]s)\n", arg.Name, arg.CType)
			}
			if f.IsVariadic {
				return strings.TrimSuffix(strings.TrimPrefix(args, "gorecv,"), ",") + "..."
			} else {
				return strings.TrimSuffix(strings.TrimPrefix(args, "gorecv,"), ",")
			}
		}
		for _, k := range keys {
			f := funcMap[k]
			if f.Return.CType != "" {
				writeOutputTypeConversion(typeMap, headerOut, f.Return.CType)
			}
		}
		for _, k := range keys {
			f := funcMap[k]
			// checkpoint, _ := headerOut.Seek(0, io.SeekCurrent)
			fmt.Fprintf(headerOut, "// export %s\n", k)
			fmt.Fprintf(headerOut, "func %s(", k)
			if len(f.Args) > 0 {
				for _, arg := range f.Args {
					if strings.HasSuffix(arg.CType, "*") {
						fmt.Fprintf(headerOut, "c%s *C.%s,", arg.Name, strings.TrimSuffix(arg.CType, "*"))
					} else {
						fmt.Fprintf(headerOut, "c%s C.%s,", arg.Name, arg.CType)
					}
				}
				headerOut.Seek(-1, io.SeekCurrent)
			}
			fmt.Fprintf(headerOut, ") ")
			if f.Return.CType != "" {
				fmt.Fprintf(headerOut, "C.%s {\n", f.Return.CType)
				if f.Recv != "" {
					fmt.Fprintf(headerOut, "\tret := gorecv.%s(%s)\n", f.TargetName, convertArgs(f))
				} else {
					fmt.Fprintf(headerOut, "\tret := %s(%s)\n", f.TargetName, convertArgs(f))
				}
				if typeMap[f.Return.CType] != nil {
					if typeMap[f.Return.CType].IsHandle {
						fmt.Fprintf(headerOut, "\treturn (C.%s)(unsafe.Pointer(cgo.NewHandle(ret)))\n", f.Return.CType)
					} else {
						fmt.Fprintf(headerOut, "\treturn convert_to_%s(ret)\n", f.Return.CType)
					}
				} else if f.Return.GoType == "bool" {
					fmt.Fprintf(headerOut, "\treturn convert_to_Go_Boolean32(ret)\n")
				} else {
					fmt.Fprintf(headerOut, "\treturn (C.%s)(ret)\n", f.Return.CType)
				}
			} else {
				fmt.Fprintf(headerOut, "{\n")
				if f.Recv != "" {
					fmt.Fprintf(headerOut, "\tgorecv.%s(%s)\n", f.TargetName, convertArgs(f))
				} else {
					fmt.Fprintf(headerOut, "\t%s(%s)\n", f.TargetName, convertArgs(f))
				}
			}
			fmt.Fprintf(headerOut, "}\n")

			/*
				if !complete {
					headerOut.Seek(checkpoint, io.SeekStart)
					headerOut.Truncate(checkpoint)
				}
			*/
		}
	}
}
