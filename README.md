# Gother

Gother gets datas gather, into one binary.

### Data Segment  
- Header size: 4 byte (32bit)
- Amount of files: 4 byte (32bit)
- All files info
    - File name: variable length
	- Code point `\0` i.e. NULL: 1 byte (8bit)
	- Data origin from .gzr file initial: 4 byte (64bit)
	- Data size: 4 byte (64bit)
- Binaries: variable length
