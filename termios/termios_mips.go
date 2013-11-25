// +build linux,arm


package termios

import (
    "os";
    "syscall";
    "unsafe"
    "errors"
)

// termios types
type cc_t byte
type speed_t uint32
type tcflag_t uint32

// termios constants
const (
    IGNBRK = tcflag_t (0000001)
    BRKINT = tcflag_t (0000002)
    IGNPAR = tcflag_t (0000004)
    PARMRK = tcflag_t (0000010)
    INLCR = tcflag_t (0000100)
    ECHONL = tcflag_t (0000100)
    IGNCR = tcflag_t (0000200)
    ICRNL = tcflag_t (0000400)
    INPCK = tcflag_t (0000020)
    ISTRIP = tcflag_t (0000040)
    IXON = tcflag_t (0002000)
    OPOST = tcflag_t (0000001)
    CS8 = tcflag_t (0000060)
    ECHO = tcflag_t (0000010)
    ICANON = tcflag_t (0000002)
    IEXTEN = tcflag_t (0100000)
    ISIG = tcflag_t (0000001)
    VTIME = tcflag_t (5)
    VMIN = tcflag_t (6)
    CBAUD = tcflag_t (0010017)
    CBAUDEX = tcflag_t (0010000)
)

const (
    B0 = speed_t(0000000)         /* hang up */
    B50 = speed_t(0000001)
    B75 = speed_t(0000002)
    B110 = speed_t(0000003)
    B134 = speed_t(0000004)
    B150 = speed_t(0000005)
    B200 = speed_t(0000006)
    B300 = speed_t(0000007)
    B600 = speed_t(0000010)
    B1200 = speed_t(0000011)
    B1800 = speed_t(0000012)
    B2400 = speed_t(0000013)
    B4800 = speed_t(0000014)
    B9600 = speed_t(0000015)
    B19200 = speed_t(0000016)
    B38400 = speed_t(0000017)
    B57600 = speed_t(0010001)
    B115200 = speed_t(0010002)
    B230400 = speed_t(0010003)
    B460800 = speed_t(0010004)
    B500000 = speed_t(0010005)
    B576000 = speed_t(0010006)
    B921600 = speed_t(0010007)
    B1000000 = speed_t(0010010)
    B1152000 = speed_t(0010011)
    B1500000 = speed_t(0010012)
    B2000000 = speed_t(0010013)
    B2500000 = speed_t(0010014)
    B3000000 = speed_t(0010015)
    B3500000 = speed_t(0010016)
    B4000000 = speed_t(0010017)
)

//note that struct termios and struct __kernel_termios have DIFFERENT size and layout !!!
const NCCS = 23  //23 on mips, 19 on alpha (also line and cc reversed), 19 on powerpc (also line and cc reversed), 17 on sparc, 
type termios struct {
    c_iflag, c_oflag, c_cflag, c_lflag  tcflag_t
    c_line  cc_t
    c_cc    [NCCS]cc_t
    c_ispeed, c_ospeed  speed_t
}

// ioctl constants
const (
    TCGETS = 0x5401
    TCSETS = 0x5402
)

func getTermios(ttyfd uintptr, dst *termios) error {
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

func setTermios(ttyfd uintptr, src *termios) error {
    r1, _, errno := syscall.Syscall (syscall.SYS_IOCTL,
                                     uintptr (ttyfd), uintptr (TCSETS),
                                     uintptr (unsafe.Pointer (src)));

    if err := os.NewSyscallError ("SYS_IOCTL", errno); errno!=0 &&err != nil {
        return err
    }

    if r1 != 0 {
        return errors.New ("Error during ioctl tcsets syscall")
    }
    return nil
}

func SetRawFd(fd uintptr) (error) {
    var orig_termios termios;
    if err := getTermios (fd, &orig_termios); err != nil { return err}

    orig_termios.c_iflag &= ^(IGNBRK|BRKINT|PARMRK|ISTRIP|INLCR|IGNCR|ICRNL|IXON);
    orig_termios.c_oflag &= ^(OPOST);
    orig_termios.c_lflag &= ^(ECHO | ECHONL | ICANON | IEXTEN | ISIG);
    orig_termios.c_cflag |= (CS8);

    orig_termios.c_cc[VMIN] = 1;
    orig_termios.c_cc[VTIME] = 0;

    return setTermios(fd, &orig_termios)
}

func SetRawFile(f *os.File) (error) {
    return SetRawFd(f.Fd())
}

func SetSpeedFd(fd uintptr, speed speed_t) (err error) {
    var orig_termios termios;
    if err = getTermios (fd, &orig_termios); err != nil { return }
    
    orig_termios.c_ispeed = speed
    orig_termios.c_ospeed = speed
    //~ //input baudrate == output baudrate and we ignore special case B0
    //~ orig_termios.c_cflag &= ^(CBAUD | CBAUDEX)
    //~ orig_termios.c_cflag |= speed
    if err = setTermios(fd, &orig_termios); err != nil { return }
    if err = getTermios (fd, &orig_termios); err != nil { return }
    if orig_termios.c_ispeed != speed || orig_termios.c_ospeed != speed {
        err = errors.New("Failed to set speed")
    }
    //~ if err = getTermios (fd, &orig_termios); err != nil { return }
    //~ if orig_termios.c_cflag & (CBAUD | CBAUDEX) != speed {
        //~ err = errors.New("Failed to set speed")
    //~ }
    return
}

func SetSpeedFile(f *os.File, speed speed_t) (error) {
    return SetSpeedFd(f.Fd(), speed)
}