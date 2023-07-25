from unittest import TestCase
from dataclasses import dataclass
from btree import varint


@dataclass
class VarIntTestCase:
    data: bytes
    expected: int
    num_bytes: int


class TestVarint(TestCase):
    def test_varint(self):
        tests = [
            VarIntTestCase([0x00], 0x00000000, 1),
            VarIntTestCase([0x7f], 0x0000007f, 1),
            VarIntTestCase([0x81, 0x00], 0x00000080, 2),
            VarIntTestCase([0x80, 0x7f], 0x0000007f, 2),
            VarIntTestCase([0x8a, 0x91, 0xd1, 0xac, 0x78], 0x12345678, 5),
            VarIntTestCase([0x81, 0x81, 0x81, 0x81, 0x01], 0x10204081, 5),
        ]
        for test in tests:
            with self.subTest(msg=f'testing varint {test.data} = {test.expected}'):
                result, num_bytes = varint(test.data, 0)
                self.assertEqual(test.expected, result)
                self.assertEqual(test.num_bytes, num_bytes)

    def test_varint_mid_sequence(self):
        raise Exception('not done yet')

    def test_varint_doesnt_read_till_end(self):
        raise Exception('not done yet')

    def test_varint_full_9_bytes(self):
        raise Exception('not done yet')

if __name__ == '__main__':
    unittest.main()
