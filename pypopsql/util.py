from typing import Tuple

def b2i(b: bytes) -> int:
    return int.from_bytes(b, 'big', signed=False)

def varint(b: bytes, cursor: int) -> Tuple[int, int]:
    """
    varint reads a variable length int from some byte string
    it takes the following parameters
     - b: sequence of bytes, generally a full table page
     - cursor: starting index of the varint

    and returns a tuple including the value of the integer, and the
    index of the cursor after the last byte of the varint
    """

    result = 0

    for j in range(8):
        i = cursor + j

        # read the next byte into an integer
        byte_num = b[i]

        # result shifts left 7 bits, then the first 7 bits of byte_num are appended
        result = (result << 7) | (byte_num & 0x7f)

        # check the first bit of byte_num to see if we should continue reading
        continue_reading = byte_num & 0x80

        if not continue_reading:
            return result, j+1

    # read last byte, use all 8 bytes to fill the remaining spaces
    byte_num = b[cursor + 8]
    result = (result << 8) | byte_num

    return result, 9

