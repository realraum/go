package termios

import (
    "fmt";
    "os";
    "syscall";
    "unsafe"
    "errors"
)

// termios types
type cc_t byte
type speed_t uint
type tcflag_t uint

// termios constants
const (
    BRKINT = tcflag_t (0000002);
    ICRNL = tcflag_t (0000400);
    INPCK = tcflag_t (0000020);
    ISTRIP = tcflag_t (0000040);
    IXON = tcflag_t (0002000);
    OPOST = tcflag_t (0000001);
    CS8 = tcflag_t (0000060);
    ECHO = tcflag_t (0000010);
    ICANON = tcflag_t (0000002);
    IEXTEN = tcflag_t (0100000);
    ISIG = tcflag_t (0000001);
    VTIME = tcflag_t (5);
    VMIN = tcflag_t (6)
)

const NCCS = 32
type termios struct {
    c_iflag, c_oflag, c_cflag, c_lflag tcflag_t;
    c_line cc_t;
    c_cc [NCCS]cc_t;
    c_ispeed, c_ospeed speed_t
}

// ioctl constants
const (
    TCGETS = 0x5401;
    TCSETS = 0x5402
)

var (
    orig_termios termios;
    ttyfd uintptr = 0 // STDIN_FILENO
)

func Ttyfd(fd uintptr) {
  ttyfd=fd
}

func getTermios (dst *termios) error {
    r1, _, errno := syscall.Syscall (syscall.SYS_IOCTL,
                                     uintptr (ttyfd), uintptr (TCGETS),
                                     uintptr (unsafe.Pointer (dst)));

    if err := os.NewSyscallError ("SYS_IOCTL", errno); errno!=0 && err != nil {
        return err
    }

    if r1 != 0 {
    //    return errors.New("Error")
    }

    return nil
}

func setTermios (src *termios) error {
    r1, _, errno := syscall.Syscall (syscall.SYS_IOCTL,
                                     uintptr (ttyfd), uintptr (TCSETS),
                                     uintptr (unsafe.Pointer (src)));

    if err := os.NewSyscallError ("SYS_IOCTL", errno); errno!=0 &&err != nil {
        return err
    }

    if r1 != 0 {
        return errors.New ("Error")
    }

    return nil
}

func tty_raw () error {
    raw := orig_termios;

    raw.c_iflag &= ^(BRKINT | ICRNL | INPCK | ISTRIP | IXON);
    raw.c_oflag &= ^(OPOST);
    raw.c_cflag |= (CS8);
    raw.c_lflag &= ^(ECHO | ICANON | IEXTEN | ISIG);

    raw.c_cc[VMIN] = 1;
    raw.c_cc[VTIME] = 0;

    if err := setTermios (&raw); err != nil { return err }

    return nil
}

func SetRaw () {
    var (
        err error
    )

    defer func () {
        if err != nil { fmt.Printf ("SetRaw Error: %v\n",err) }
    } ();

    if err = getTermios (&orig_termios); err != nil { return }

//    defer func () {
//        err = setTermios (&orig_termios)
//    } ();

    if err = tty_raw (); err != nil { return }
    //if err = screenio (); err != nil { return }
}

func SetSpeed (speed uint) {
    var err error

    defer func () {
        if err != nil { fmt.Printf ("SetSpeed Error: %v\n",err) }
    } ();

    if err = getTermios (&orig_termios); err != nil { return }
    orig_termios.c_ispeed = speed_t(speed)
    orig_termios.c_ospeed = speed_t(speed)
    err = setTermios (&orig_termios)
}
