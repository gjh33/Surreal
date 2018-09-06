package input

import "github.com/go-gl/glfw/v3.2/glfw"

// KeyCode enum to represent keyboard keys
type KeyCode int

// KeyCode value declarations
const (
	KeyUnknown      KeyCode = KeyCode(glfw.KeyUnknown)
	KeySpace        KeyCode = KeyCode(glfw.KeySpace)
	KeyApostrophe   KeyCode = KeyCode(glfw.KeyApostrophe)
	KeyComma        KeyCode = KeyCode(glfw.KeyComma)
	KeyMinus        KeyCode = KeyCode(glfw.KeyMinus)
	KeyPeriod       KeyCode = KeyCode(glfw.KeyPeriod)
	KeySlash        KeyCode = KeyCode(glfw.KeySlash)
	Key0            KeyCode = KeyCode(glfw.Key0)
	Key1            KeyCode = KeyCode(glfw.Key1)
	Key2            KeyCode = KeyCode(glfw.Key2)
	Key3            KeyCode = KeyCode(glfw.Key3)
	Key4            KeyCode = KeyCode(glfw.Key4)
	Key5            KeyCode = KeyCode(glfw.Key5)
	Key6            KeyCode = KeyCode(glfw.Key6)
	Key7            KeyCode = KeyCode(glfw.Key7)
	Key8            KeyCode = KeyCode(glfw.Key8)
	Key9            KeyCode = KeyCode(glfw.Key9)
	KeySemicolon    KeyCode = KeyCode(glfw.KeySemicolon)
	KeyEqual        KeyCode = KeyCode(glfw.KeyEqual)
	KeyA            KeyCode = KeyCode(glfw.KeyA)
	KeyB            KeyCode = KeyCode(glfw.KeyB)
	KeyC            KeyCode = KeyCode(glfw.KeyC)
	KeyD            KeyCode = KeyCode(glfw.KeyD)
	KeyE            KeyCode = KeyCode(glfw.KeyE)
	KeyF            KeyCode = KeyCode(glfw.KeyF)
	KeyG            KeyCode = KeyCode(glfw.KeyG)
	KeyH            KeyCode = KeyCode(glfw.KeyH)
	KeyI            KeyCode = KeyCode(glfw.KeyI)
	KeyJ            KeyCode = KeyCode(glfw.KeyJ)
	KeyK            KeyCode = KeyCode(glfw.KeyK)
	KeyL            KeyCode = KeyCode(glfw.KeyL)
	KeyM            KeyCode = KeyCode(glfw.KeyM)
	KeyN            KeyCode = KeyCode(glfw.KeyN)
	KeyO            KeyCode = KeyCode(glfw.KeyO)
	KeyP            KeyCode = KeyCode(glfw.KeyP)
	KeyQ            KeyCode = KeyCode(glfw.KeyQ)
	KeyR            KeyCode = KeyCode(glfw.KeyR)
	KeyS            KeyCode = KeyCode(glfw.KeyS)
	KeyT            KeyCode = KeyCode(glfw.KeyT)
	KeyU            KeyCode = KeyCode(glfw.KeyU)
	KeyV            KeyCode = KeyCode(glfw.KeyV)
	KeyW            KeyCode = KeyCode(glfw.KeyW)
	KeyX            KeyCode = KeyCode(glfw.KeyX)
	KeyY            KeyCode = KeyCode(glfw.KeyY)
	KeyZ            KeyCode = KeyCode(glfw.KeyZ)
	KeyLeftBracket  KeyCode = KeyCode(glfw.KeyLeftBracket)
	KeyBackslash    KeyCode = KeyCode(glfw.KeyBackslash)
	KeyRightBracket KeyCode = KeyCode(glfw.KeyRightBracket)
	KeyGraveAccent  KeyCode = KeyCode(glfw.KeyGraveAccent)
	KeyWorld1       KeyCode = KeyCode(glfw.KeyWorld1)
	KeyWorld2       KeyCode = KeyCode(glfw.KeyWorld2)
	KeyEscape       KeyCode = KeyCode(glfw.KeyEscape)
	KeyEnter        KeyCode = KeyCode(glfw.KeyEnter)
	KeyTab          KeyCode = KeyCode(glfw.KeyTab)
	KeyBackspace    KeyCode = KeyCode(glfw.KeyBackspace)
	KeyInsert       KeyCode = KeyCode(glfw.KeyInsert)
	KeyDelete       KeyCode = KeyCode(glfw.KeyDelete)
	KeyRight        KeyCode = KeyCode(glfw.KeyRight)
	KeyLeft         KeyCode = KeyCode(glfw.KeyLeft)
	KeyDown         KeyCode = KeyCode(glfw.KeyDown)
	KeyUp           KeyCode = KeyCode(glfw.KeyUp)
	KeyPageUp       KeyCode = KeyCode(glfw.KeyPageUp)
	KeyPageDown     KeyCode = KeyCode(glfw.KeyPageDown)
	KeyHome         KeyCode = KeyCode(glfw.KeyHome)
	KeyEnd          KeyCode = KeyCode(glfw.KeyEnd)
	KeyCapsLock     KeyCode = KeyCode(glfw.KeyCapsLock)
	KeyScrollLock   KeyCode = KeyCode(glfw.KeyScrollLock)
	KeyNumLock      KeyCode = KeyCode(glfw.KeyNumLock)
	KeyPrintScreen  KeyCode = KeyCode(glfw.KeyPrintScreen)
	KeyPause        KeyCode = KeyCode(glfw.KeyPause)
	KeyF1           KeyCode = KeyCode(glfw.KeyF1)
	KeyF2           KeyCode = KeyCode(glfw.KeyF2)
	KeyF3           KeyCode = KeyCode(glfw.KeyF3)
	KeyF4           KeyCode = KeyCode(glfw.KeyF4)
	KeyF5           KeyCode = KeyCode(glfw.KeyF5)
	KeyF6           KeyCode = KeyCode(glfw.KeyF6)
	KeyF7           KeyCode = KeyCode(glfw.KeyF7)
	KeyF8           KeyCode = KeyCode(glfw.KeyF8)
	KeyF9           KeyCode = KeyCode(glfw.KeyF9)
	KeyF10          KeyCode = KeyCode(glfw.KeyF10)
	KeyF11          KeyCode = KeyCode(glfw.KeyF11)
	KeyF12          KeyCode = KeyCode(glfw.KeyF12)
	KeyF13          KeyCode = KeyCode(glfw.KeyF13)
	KeyF14          KeyCode = KeyCode(glfw.KeyF14)
	KeyF15          KeyCode = KeyCode(glfw.KeyF15)
	KeyF16          KeyCode = KeyCode(glfw.KeyF16)
	KeyF17          KeyCode = KeyCode(glfw.KeyF17)
	KeyF18          KeyCode = KeyCode(glfw.KeyF18)
	KeyF19          KeyCode = KeyCode(glfw.KeyF19)
	KeyF20          KeyCode = KeyCode(glfw.KeyF20)
	KeyF21          KeyCode = KeyCode(glfw.KeyF21)
	KeyF22          KeyCode = KeyCode(glfw.KeyF22)
	KeyF23          KeyCode = KeyCode(glfw.KeyF23)
	KeyF24          KeyCode = KeyCode(glfw.KeyF24)
	KeyF25          KeyCode = KeyCode(glfw.KeyF25)
	KeyKP0          KeyCode = KeyCode(glfw.KeyKP0)
	KeyKP1          KeyCode = KeyCode(glfw.KeyKP1)
	KeyKP2          KeyCode = KeyCode(glfw.KeyKP2)
	KeyKP3          KeyCode = KeyCode(glfw.KeyKP3)
	KeyKP4          KeyCode = KeyCode(glfw.KeyKP4)
	KeyKP5          KeyCode = KeyCode(glfw.KeyKP5)
	KeyKP6          KeyCode = KeyCode(glfw.KeyKP6)
	KeyKP7          KeyCode = KeyCode(glfw.KeyKP7)
	KeyKP8          KeyCode = KeyCode(glfw.KeyKP8)
	KeyKP9          KeyCode = KeyCode(glfw.KeyKP9)
	KeyKPDecimal    KeyCode = KeyCode(glfw.KeyKPDecimal)
	KeyKPDivide     KeyCode = KeyCode(glfw.KeyKPDivide)
	KeyKPMultiply   KeyCode = KeyCode(glfw.KeyKPMultiply)
	KeyKPSubtract   KeyCode = KeyCode(glfw.KeyKPSubtract)
	KeyKPAdd        KeyCode = KeyCode(glfw.KeyKPAdd)
	KeyKPEnter      KeyCode = KeyCode(glfw.KeyKPEnter)
	KeyKPEqual      KeyCode = KeyCode(glfw.KeyKPEqual)
	KeyLeftShift    KeyCode = KeyCode(glfw.KeyLeftShift)
	KeyLeftControl  KeyCode = KeyCode(glfw.KeyLeftControl)
	KeyLeftAlt      KeyCode = KeyCode(glfw.KeyLeftAlt)
	KeyLeftSuper    KeyCode = KeyCode(glfw.KeyLeftSuper)
	KeyRightShift   KeyCode = KeyCode(glfw.KeyRightShift)
	KeyRightControl KeyCode = KeyCode(glfw.KeyRightControl)
	KeyRightAlt     KeyCode = KeyCode(glfw.KeyRightAlt)
	KeyRightSuper   KeyCode = KeyCode(glfw.KeyRightSuper)
	KeyMenu         KeyCode = KeyCode(glfw.KeyMenu)
	KeyLast         KeyCode = KeyCode(glfw.KeyLast)
)

// KeyState represents the key's current state
type KeyState int

// The enum values for KeyState
const (
	PressedState  KeyState = KeyState(glfw.Press)
	ReleasedState KeyState = KeyState(glfw.Release)
	StayState     KeyState = KeyState(glfw.Repeat)
)

var state = make(map[KeyCode]KeyState)

// TrackWindow tells the input system to start tracking input from this window
func TrackWindow(win *glfw.Window) {
	win.SetKeyCallback(glfwKeyboardCallback)
}

// GetKeyDown takes a keycode and returns true if the key was pressed this frame
func GetKeyDown(key KeyCode) bool {
	val, ok := state[key]
	if !ok {
		return false
	}
	return val == PressedState
}

// GetKeyUp takes a keycode and returns true if the key was released this frame
func GetKeyUp(key KeyCode) bool {
	val, ok := state[key]
	if !ok {
		return false
	}
	return val == ReleasedState
}

// GetKey takes a keycode and returns true if the key is currently down this frame including if it was just pressed
func GetKey(key KeyCode) bool {
	val, ok := state[key]
	if !ok {
		return false
	}
	return val == StayState || val == PressedState
}

func glfwKeyboardCallback(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	state[KeyCode(key)] = KeyState(action)
}
