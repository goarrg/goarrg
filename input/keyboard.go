/*
Copyright 2020 The goARRG Authors.

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

package input

/*
Key codes of keys on the keyboard as defined in the USB HID Usage Tables
https://usb.org/document-library/hid-usage-tables-122
*/
const (
	KeyA                   DeviceAction = 4
	KeyB                   DeviceAction = 5
	KeyC                   DeviceAction = 6
	KeyD                   DeviceAction = 7
	KeyE                   DeviceAction = 8
	KeyF                   DeviceAction = 9
	KeyG                   DeviceAction = 10
	KeyH                   DeviceAction = 11
	KeyI                   DeviceAction = 12
	KeyJ                   DeviceAction = 13
	KeyK                   DeviceAction = 14
	KeyL                   DeviceAction = 15
	KeyM                   DeviceAction = 16
	KeyN                   DeviceAction = 17
	KeyO                   DeviceAction = 18
	KeyP                   DeviceAction = 19
	KeyQ                   DeviceAction = 20
	KeyR                   DeviceAction = 21
	KeyS                   DeviceAction = 22
	KeyT                   DeviceAction = 23
	KeyU                   DeviceAction = 24
	KeyV                   DeviceAction = 25
	KeyW                   DeviceAction = 26
	KeyX                   DeviceAction = 27
	KeyY                   DeviceAction = 28
	KeyZ                   DeviceAction = 29
	Key1                   DeviceAction = 30
	Key2                   DeviceAction = 31
	Key3                   DeviceAction = 32
	Key4                   DeviceAction = 33
	Key5                   DeviceAction = 34
	Key6                   DeviceAction = 35
	Key7                   DeviceAction = 36
	Key8                   DeviceAction = 37
	Key9                   DeviceAction = 38
	Key0                   DeviceAction = 39
	KeyEnter               DeviceAction = 40
	KeyEscape              DeviceAction = 41
	KeyBackspace           DeviceAction = 42
	KeyTab                 DeviceAction = 43
	KeySpacebar            DeviceAction = 44
	KeyMinus               DeviceAction = 45
	KeyEquals              DeviceAction = 46
	KeyLeftBracket         DeviceAction = 47
	KeyRightBracket        DeviceAction = 48
	KeyBackslash           DeviceAction = 49
	KeyNonUSHash           DeviceAction = 50
	KeySemicolon           DeviceAction = 51
	KeyApostrophe          DeviceAction = 52
	KeyGrave               DeviceAction = 53
	KeyComma               DeviceAction = 54
	KeyPeriod              DeviceAction = 55
	KeySlash               DeviceAction = 56
	KeyCapsLock            DeviceAction = 57
	KeyF1                  DeviceAction = 58
	KeyF2                  DeviceAction = 59
	KeyF3                  DeviceAction = 60
	KeyF4                  DeviceAction = 61
	KeyF5                  DeviceAction = 62
	KeyF6                  DeviceAction = 63
	KeyF7                  DeviceAction = 64
	KeyF8                  DeviceAction = 65
	KeyF9                  DeviceAction = 66
	KeyF10                 DeviceAction = 67
	KeyF11                 DeviceAction = 68
	KeyF12                 DeviceAction = 69
	KeyPrintScreen         DeviceAction = 70
	KeyScrollLock          DeviceAction = 71
	KeyPause               DeviceAction = 72
	KeyInsert              DeviceAction = 73
	KeyHome                DeviceAction = 74
	KeyPageUp              DeviceAction = 75
	KeyDelete              DeviceAction = 76
	KeyEnd                 DeviceAction = 77
	KeyPageDown            DeviceAction = 78
	KeyRight               DeviceAction = 79
	KeyLeft                DeviceAction = 80
	KeyDown                DeviceAction = 81
	KeyUp                  DeviceAction = 82
	KeyNumLock             DeviceAction = 83
	KeyKPDivide            DeviceAction = 84
	KeyKPMultiply          DeviceAction = 85
	KeyKPMinus             DeviceAction = 86
	KeyKPPlus              DeviceAction = 87
	KeyKPEnter             DeviceAction = 88
	KeyKP1                 DeviceAction = 89
	KeyKP2                 DeviceAction = 90
	KeyKP3                 DeviceAction = 91
	KeyKP4                 DeviceAction = 92
	KeyKP5                 DeviceAction = 93
	KeyKP6                 DeviceAction = 94
	KeyKP7                 DeviceAction = 95
	KeyKP8                 DeviceAction = 96
	KeyKP9                 DeviceAction = 97
	KeyKP0                 DeviceAction = 98
	KeyKPPeriod            DeviceAction = 99
	KeyNonUSBackslash      DeviceAction = 100
	KeyApplication         DeviceAction = 101
	KeyPower               DeviceAction = 102
	KeyKPEquals            DeviceAction = 103
	KeyF13                 DeviceAction = 104
	KeyF14                 DeviceAction = 105
	KeyF15                 DeviceAction = 106
	KeyF16                 DeviceAction = 107
	KeyF17                 DeviceAction = 108
	KeyF18                 DeviceAction = 109
	KeyF19                 DeviceAction = 110
	KeyF20                 DeviceAction = 111
	KeyF21                 DeviceAction = 112
	KeyF22                 DeviceAction = 113
	KeyF23                 DeviceAction = 114
	KeyF24                 DeviceAction = 115
	KeyExecute             DeviceAction = 116
	KeyHelp                DeviceAction = 117
	KeyMenu                DeviceAction = 118
	KeySelect              DeviceAction = 119
	KeyStop                DeviceAction = 120
	KeyAgain               DeviceAction = 121
	KeyUndo                DeviceAction = 122
	KeyCut                 DeviceAction = 123
	KeyCopy                DeviceAction = 124
	KeyPaste               DeviceAction = 125
	KeyFind                DeviceAction = 126
	KeyMute                DeviceAction = 127
	KeyVolumeUp            DeviceAction = 128
	KeyVolumeDown          DeviceAction = 129
	KeyKPComma             DeviceAction = 133
	KeyKPEqualsAS400       DeviceAction = 134
	KeyInternational1      DeviceAction = 135
	KeyInternational2      DeviceAction = 136
	KeyInternational3      DeviceAction = 137
	KeyInternational4      DeviceAction = 138
	KeyInternational5      DeviceAction = 139
	KeyInternational6      DeviceAction = 140
	KeyInternational7      DeviceAction = 141
	KeyInternational8      DeviceAction = 142
	KeyInternational9      DeviceAction = 143
	KeyLang1               DeviceAction = 144
	KeyLang2               DeviceAction = 145
	KeyLang3               DeviceAction = 146
	KeyLang4               DeviceAction = 147
	KeyLang5               DeviceAction = 148
	KeyLang6               DeviceAction = 149
	KeyLang7               DeviceAction = 150
	KeyLang8               DeviceAction = 151
	KeyLang9               DeviceAction = 152
	KeyAltErase            DeviceAction = 153
	KeySysReq              DeviceAction = 154
	KeyCancel              DeviceAction = 155
	KeyClear               DeviceAction = 156
	KeyPrior               DeviceAction = 157
	KeyReturn              DeviceAction = 158
	KeySeparator           DeviceAction = 159
	KeyOut                 DeviceAction = 160
	KeyOper                DeviceAction = 161
	KeyClearAgain          DeviceAction = 162
	KeyCrSel               DeviceAction = 163
	KeyExSel               DeviceAction = 164
	KeyKP00                DeviceAction = 176
	KeyKP000               DeviceAction = 177
	KeyThousandsSeparator  DeviceAction = 178
	KeyDecimalSeparator    DeviceAction = 179
	KeyCurrencyUnit        DeviceAction = 180
	KeyCurrencySubmit      DeviceAction = 181
	KeyKPLeftParentheses   DeviceAction = 182
	KeyKPRightParentheses  DeviceAction = 183
	KeyKPLeftBrace         DeviceAction = 184
	KeyKPRightBrace        DeviceAction = 185
	KeyKPTab               DeviceAction = 186
	KeyKPBackspace         DeviceAction = 187
	KeyKPA                 DeviceAction = 188
	KeyKPB                 DeviceAction = 189
	KeyKPC                 DeviceAction = 190
	KeyKPD                 DeviceAction = 191
	KeyKPE                 DeviceAction = 192
	KeyKPF                 DeviceAction = 193
	KeyKPXOR               DeviceAction = 194
	KeyKPPower             DeviceAction = 195
	KeyKPPercent           DeviceAction = 196
	KeyKPLess              DeviceAction = 197
	KeyKPGreater           DeviceAction = 198
	KeyKPAmpersand         DeviceAction = 199
	KeyKPDoubleAmpersand   DeviceAction = 200
	KeyKPVerticalBar       DeviceAction = 201
	KeyKPDoubleVerticalBar DeviceAction = 202
	KeyKPColon             DeviceAction = 203
	KeyKPHash              DeviceAction = 204
	KeyKPSpace             DeviceAction = 205
	KeyKPAt                DeviceAction = 206
	KeyKPExclamation       DeviceAction = 207
	KeyKPMemStore          DeviceAction = 208
	KeyKPMemRecall         DeviceAction = 209
	KeyKPMemClear          DeviceAction = 210
	KeyKPMemAdd            DeviceAction = 211
	KeyKPMemSubtract       DeviceAction = 212
	KeyKPMemMultiply       DeviceAction = 213
	KeyKPMemDivide         DeviceAction = 214
	KeyKPPlusMinus         DeviceAction = 215
	KeyKPClear             DeviceAction = 216
	KeyKPClearEntry        DeviceAction = 217
	KeyKPBinary            DeviceAction = 218
	KeyKPOctal             DeviceAction = 219
	KeyKPDecimal           DeviceAction = 220
	KeyKPHexadecimal       DeviceAction = 221
	KeyLeftCtrl            DeviceAction = 224
	KeyLeftShift           DeviceAction = 225
	KeyLeftAlt             DeviceAction = 226
	KeyLeftGUI             DeviceAction = 227
	KeyRightCtrl           DeviceAction = 228
	KeyRightShift          DeviceAction = 229
	KeyRightAlt            DeviceAction = 230
	KeyRightGUI            DeviceAction = 231

	// Deprecated keys
	// KeyLockingCapsLock     DeviceAction = 130
	// KeyLockingNumLock      DeviceAction = 131
	// KeyLockingScrollLock   DeviceAction = 132
)

var _ = map[DeviceAction]string{
	KeyA:                   "A",
	KeyB:                   "B",
	KeyC:                   "C",
	KeyD:                   "D",
	KeyE:                   "E",
	KeyF:                   "F",
	KeyG:                   "G",
	KeyH:                   "H",
	KeyI:                   "I",
	KeyJ:                   "J",
	KeyK:                   "K",
	KeyL:                   "L",
	KeyM:                   "M",
	KeyN:                   "N",
	KeyO:                   "O",
	KeyP:                   "P",
	KeyQ:                   "Q",
	KeyR:                   "R",
	KeyS:                   "S",
	KeyT:                   "T",
	KeyU:                   "U",
	KeyV:                   "V",
	KeyW:                   "W",
	KeyX:                   "X",
	KeyY:                   "Y",
	KeyZ:                   "Z",
	Key1:                   "1",
	Key2:                   "2",
	Key3:                   "3",
	Key4:                   "4",
	Key5:                   "5",
	Key6:                   "6",
	Key7:                   "7",
	Key8:                   "8",
	Key9:                   "9",
	Key0:                   "0",
	KeyEnter:               "Enter",
	KeyEscape:              "Escape",
	KeyBackspace:           "Backspace",
	KeyTab:                 "Tab",
	KeySpacebar:            "Spacebar",
	KeyMinus:               "Minus",
	KeyEquals:              "Equals",
	KeyLeftBracket:         "LeftBracket",
	KeyRightBracket:        "RightBracket",
	KeyBackslash:           "Backslash",
	KeyNonUSHash:           "NonUSHash",
	KeySemicolon:           "Semicolon",
	KeyApostrophe:          "Apostrophe",
	KeyGrave:               "Grave",
	KeyComma:               "Comma",
	KeyPeriod:              "Period",
	KeySlash:               "Slash",
	KeyCapsLock:            "CapsLock",
	KeyF1:                  "F1",
	KeyF2:                  "F2",
	KeyF3:                  "F3",
	KeyF4:                  "F4",
	KeyF5:                  "F5",
	KeyF6:                  "F6",
	KeyF7:                  "F7",
	KeyF8:                  "F8",
	KeyF9:                  "F9",
	KeyF10:                 "F10",
	KeyF11:                 "F11",
	KeyF12:                 "F12",
	KeyPrintScreen:         "PrintScreen",
	KeyScrollLock:          "ScrollLock",
	KeyPause:               "Pause",
	KeyInsert:              "Insert",
	KeyHome:                "Home",
	KeyPageUp:              "PageUp",
	KeyDelete:              "Delete",
	KeyEnd:                 "End",
	KeyPageDown:            "PageDown",
	KeyRight:               "Right",
	KeyLeft:                "Left",
	KeyDown:                "Down",
	KeyUp:                  "Up",
	KeyNumLock:             "NumLock",
	KeyKPDivide:            "KPDivide",
	KeyKPMultiply:          "KPMultiply",
	KeyKPMinus:             "KPMinus",
	KeyKPPlus:              "KPPlus",
	KeyKPEnter:             "KPEnter",
	KeyKP1:                 "KP1",
	KeyKP2:                 "KP2",
	KeyKP3:                 "KP3",
	KeyKP4:                 "KP4",
	KeyKP5:                 "KP5",
	KeyKP6:                 "KP6",
	KeyKP7:                 "KP7",
	KeyKP8:                 "KP8",
	KeyKP9:                 "KP9",
	KeyKP0:                 "KP0",
	KeyKPPeriod:            "KPPeriod",
	KeyNonUSBackslash:      "NonUSBackslash",
	KeyApplication:         "Application",
	KeyPower:               "Power",
	KeyKPEquals:            "KPEquals",
	KeyF13:                 "F13",
	KeyF14:                 "F14",
	KeyF15:                 "F15",
	KeyF16:                 "F16",
	KeyF17:                 "F17",
	KeyF18:                 "F18",
	KeyF19:                 "F19",
	KeyF20:                 "F20",
	KeyF21:                 "F21",
	KeyF22:                 "F22",
	KeyF23:                 "F23",
	KeyF24:                 "F24",
	KeyExecute:             "Execute",
	KeyHelp:                "Help",
	KeyMenu:                "Menu",
	KeySelect:              "Select",
	KeyStop:                "Stop",
	KeyAgain:               "Again",
	KeyUndo:                "Undo",
	KeyCut:                 "Cut",
	KeyCopy:                "Copy",
	KeyPaste:               "Paste",
	KeyFind:                "Find",
	KeyMute:                "Mute",
	KeyVolumeUp:            "VolumeUp",
	KeyVolumeDown:          "VolumeDown",
	KeyKPComma:             "KPComma",
	KeyKPEqualsAS400:       "KPEqualsAS400",
	KeyInternational1:      "International1",
	KeyInternational2:      "International2",
	KeyInternational3:      "International3",
	KeyInternational4:      "International4",
	KeyInternational5:      "International5",
	KeyInternational6:      "International6",
	KeyInternational7:      "International7",
	KeyInternational8:      "International8",
	KeyInternational9:      "International9",
	KeyLang1:               "Lang1",
	KeyLang2:               "Lang2",
	KeyLang3:               "Lang3",
	KeyLang4:               "Lang4",
	KeyLang5:               "Lang5",
	KeyLang6:               "Lang6",
	KeyLang7:               "Lang7",
	KeyLang8:               "Lang8",
	KeyLang9:               "Lang9",
	KeyAltErase:            "AltErase",
	KeySysReq:              "SysReq",
	KeyCancel:              "Cancel",
	KeyClear:               "Clear",
	KeyPrior:               "Prior",
	KeyReturn:              "Return",
	KeySeparator:           "Separator",
	KeyOut:                 "Out",
	KeyOper:                "Oper",
	KeyClearAgain:          "ClearAgain",
	KeyCrSel:               "CrSel",
	KeyExSel:               "ExSel",
	KeyKP00:                "KP00",
	KeyKP000:               "KP000",
	KeyThousandsSeparator:  "ThousandsSeparator",
	KeyDecimalSeparator:    "DecimalSeparator",
	KeyCurrencyUnit:        "CurrencyUnit",
	KeyCurrencySubmit:      "CurrencySubmit",
	KeyKPLeftParentheses:   "KPLeftParentheses",
	KeyKPRightParentheses:  "KPRightParentheses",
	KeyKPLeftBrace:         "KPLeftBrace",
	KeyKPRightBrace:        "KPRightBrace",
	KeyKPTab:               "KPTab",
	KeyKPBackspace:         "KPBackspace",
	KeyKPA:                 "KPA",
	KeyKPB:                 "KPB",
	KeyKPC:                 "KPC",
	KeyKPD:                 "KPD",
	KeyKPE:                 "KPE",
	KeyKPF:                 "KPF",
	KeyKPXOR:               "KPXOR",
	KeyKPPower:             "KPPower",
	KeyKPPercent:           "KPPercent",
	KeyKPLess:              "KPLess",
	KeyKPGreater:           "KPGreater",
	KeyKPAmpersand:         "KPAmpersand",
	KeyKPDoubleAmpersand:   "KPDoubleAmpersand",
	KeyKPVerticalBar:       "KPVerticalBar",
	KeyKPDoubleVerticalBar: "KPDoubleVerticalBar",
	KeyKPColon:             "KPColon",
	KeyKPHash:              "KPHash",
	KeyKPSpace:             "KPSpace",
	KeyKPAt:                "KPAt",
	KeyKPExclamation:       "KPExclamation",
	KeyKPMemStore:          "KPMemStore",
	KeyKPMemRecall:         "KPMemRecall",
	KeyKPMemClear:          "KPMemClear",
	KeyKPMemAdd:            "KPMemAdd",
	KeyKPMemSubtract:       "KPMemSubtract",
	KeyKPMemMultiply:       "KPMemMultiply",
	KeyKPMemDivide:         "KPMemDivide",
	KeyKPPlusMinus:         "KPPlusMinus",
	KeyKPClear:             "KPClear",
	KeyKPClearEntry:        "KPClearEntry",
	KeyKPBinary:            "KPBinary",
	KeyKPOctal:             "KPOctal",
	KeyKPDecimal:           "KPDecimal",
	KeyKPHexadecimal:       "KPHexadecimal",
	KeyLeftCtrl:            "LeftCtrl",
	KeyLeftShift:           "LeftShift",
	KeyLeftAlt:             "LeftAlt",
	KeyLeftGUI:             "LeftGUI",
	KeyRightCtrl:           "RightCtrl",
	KeyRightShift:          "RightShift",
	KeyRightAlt:            "RightAlt",
	KeyRightGUI:            "RightGUI",
}
