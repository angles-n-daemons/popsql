from header import Header, FileFormatVersion, SchemaFormat, TextEncoding, Version
from unittest import TestCase

class TestHeader(TestCase):
    def test_header_parse(self):
        data = b'SQLite format 3\x00' + bytes([
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
        
        header = Header(data)
        self.assertEqual(header.page_size, 4096)
        self.assertEqual(header.file_format_write_version, FileFormatVersion.LEGACY)
        self.assertEqual(header.file_format_read_version, FileFormatVersion.LEGACY)
        self.assertEqual(header.page_end_reserved_space, 12)
        self.assertEqual(header.maximum_embedded_payload_fraction, 64)
        self.assertEqual(header.minimum_embedded_payload_fraction, 32)
        self.assertEqual(header.leaf_payload_fraction, 32)
        self.assertEqual(header.file_change_counter, 2)
        self.assertEqual(header.db_size_in_pages, 2)
        self.assertEqual(header.first_freelist_trunk_page, 3)
        self.assertEqual(header.num_freelist_pages, 1)
        self.assertEqual(header.schema_cookie, 1)
        self.assertEqual(header.schema_format_number, SchemaFormat.FORMAT_4)
        self.assertEqual(header.default_page_cache_size, 16)
        self.assertEqual(header.largest_btree_root_page, 0)
        self.assertEqual(header.text_encoding, TextEncoding.UTF_8)
        self.assertEqual(header.user_version, 0)
        self.assertEqual(header.incremental_vacuum_mode, 0)
        self.assertEqual(header.application_id, 0)
        self.assertEqual(header.version_valid_for, 2)
        self.assertEqual(header.version, Version(3, 39, 5))

if __name__ == '__main__':
    unittest.main()
