# Region Dump Format Specification
The region dump should be a binary file named `REGION_DUMP` and stored in the build sync plugin directory.

## Endianness
All numbers in this format are stored as big-endian. Integers are 4 bytes.

## Palette
| Field | Type | Description |
|-------|------|-------------|
| Size | Integer | The number of entries in the palette |
| Palette Entries | Array of Palette Entry | The palette entries, stored side-to-side. Each palette entry is basically just a string. The specifics of how it is stored is given in the table below. |

This is how each palette entry is stored:
| Field | Type | Description |
|-------|------|-------------|
| Length | Integer | The number of characters/bytes in the palette entry |
| Name | String | The full name of the block |

## Block Data
| Field | Type | Description |
|-------|------|-------------|
| Blocks | Array of Integer | The 256x256x256 array of blocks, stored as integers. The entries are indexed by `x + z*256 + y*256*256`. |