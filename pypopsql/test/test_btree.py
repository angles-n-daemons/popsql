from btree import Node
from unittest import TestCase

EXAMPLE_HEADER = b'SQLite format 3\x00' + bytes([
    # page size
    0x10, 0x00,

    # file format write version
    0x01,
    # file format read version
    0x01,

    # page end reserved space
    0x0c,

    # maximum embedded payload fraction
    0x40,
    # minimum embedded payload fraction
    0x20,
    # leaf payload fraction
    0x20,

    # file change counter
    0x00, 0x00, 0x00, 0x02,

    # size of database in pages
    0x00, 0x00, 0x00, 0x02,

    # page number of the first freelist trunk page
    0x00, 0x00, 0x00, 0x03,
    # total number of freelist pages
    0x00, 0x00, 0x00, 0x01,

    # schema cookie
    0x00, 0x00, 0x00, 0x01,
    # schema format number
    0x00, 0x00, 0x00, 0x04,

    # default page cache size
    0x00, 0x00, 0x00, 0x10,
    # page number of the largest root b-tree page if in auto-vacuum or incremental vacuum mode
    0x00, 0x00, 0x00, 0x00,
    # text encoding
    0x00, 0x00, 0x00, 0x01,

    # user version
    0x00, 0x00, 0x00, 0x00,
    # incremental vacuum mode
    0x00, 0x00, 0x00, 0x00,
    # application id
    0x00, 0x00, 0x00, 0x00,

    # space reserved for expansion
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,

    # version-valid-for
    0x00, 0x00, 0x00, 0x02,
    # sqlite version number 3.39.5
    0x00, 0x2e, 0x5f, 0x1d,

    # random bytes to ensure it stops parsing appropriately
    0x00, 0x01, 0x02, 0x03,
    0x04, 0x05, 0x06, 0x07,
])

class TestBTree(TestCase):
    def test_table_leaf_parse(self):
        data = bytes([
            0x0d, # Node Type 13
            0x00, 0x00, # first_freeblock
            0x00, 0x02, # num cells 2
            0x00, 0x11, # cell offset 17
            0x00, # num fragmented bytes

            ## cell pointer array
            # row 1 offset 25
            0x00, 0x19,
            # row 2 offset 17
            0x00, 0x11,

            # null bytes in the center of the page
            0x00, 0x00, 0x00, 0x00, 0x00,

            # row 2
            0x06, # payload size
            0x02, # row id
            0x03, 0x11, 0x01, 0x79, 0x6f, 0x02, # payload

            # row 1
            0x05, # payload size
            0x01, # row id
            0x03, 0x11, 0x09, 0x68, 0x69, # payload
        ])

        node = Node(data)
        self.assertEqual(node.data, data)
        self.assertEqual(node.page_size, 32)
        self.assertEqual(node.num_cells, 2)
        self.assertEqual(node.cell_offset, 17)
        self.assertEqual(node.first_freeblock, 0)
        self.assertEqual(node.num_fragmented_bytes, 0)
        self.assertEqual(node.right_pointer, None)
        self.assertEqual(len(node.cells), 2)

        # check first row
        self.assertEqual(node.cells[0].payload_size, 5)
        self.assertEqual(node.cells[0].row_id, 1)
        self.assertEqual(node.cells[0].payload, bytes([0x03, 0x11, 0x09, 0x68, 0x69]))
        self.assertEqual(node.cells[0].cursor, 32)

        # check second row
        self.assertEqual(node.cells[1].payload_size, 6)
        self.assertEqual(node.cells[1].row_id, 2)
        self.assertEqual(node.cells[1].payload, bytes([0x03, 0x11, 0x01, 0x79, 0x6f, 0x02]))
        self.assertEqual(node.cells[1].cursor, 25)

    def test_schema_header_page(self):
        raise Exception('not done')
        pass

if __name__ == '__main__':
    unittest.main()
