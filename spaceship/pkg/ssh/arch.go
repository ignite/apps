package ssh

var archMap = map[string]string{
	"x86_64":   "amd64",    // 64-bit x86 architecture, common in desktop and server processors
	"aarch64":  "arm64",    // 64-bit ARM architecture, used in many mobile and embedded devices
	"i386":     "386",      // 32-bit x86 architecture, older systems
	"armv7l":   "arm",      // 32-bit ARM architecture, commonly used in mobile and embedded devices
	"armv6l":   "arm",      // ARMv6 architecture, older 32-bit ARM systems
	"mips":     "mips",     // MIPS architecture, used in embedded systems and older computing
	"mipsle":   "mipsle",   // MIPS little-endian, a variation of MIPS
	"mips64":   "mips64",   // 64-bit MIPS architecture, used in specialized computing environments
	"mips64le": "mips64le", // 64-bit MIPS little-endian, a variation of MIPS
	"ppc64":    "ppc64",    // 64-bit PowerPC architecture, used in servers and high-performance computing
	"ppc64le":  "ppc64le",  // 64-bit PowerPC little-endian, a variation of PowerPC
	"s390x":    "s390x",    // IBM System/390 architecture, used in mainframes
	"riscv64":  "riscv64",  // 64-bit RISC-V architecture, an open-source ISA
}
