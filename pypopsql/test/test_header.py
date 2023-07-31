from unittest import TestCase

class TestHeader(TestCase):
    def test_header_parse(self):
        data = b'SQLite format 3\x00' + bytes([
            # page size
            # file format write version
            # file format read version
            # page end reserved space
            # maximum embedded payload fraction
            # minimum embedded payload fraction
            # leaf payload fraction
            # file change counter
            # size of database in pages
            # page number of the first freelist trunk page
            # total number of freelist pages
            # schema cookie
            # schema format number
            # default page cache size
            # page number of the largest root b-tree page if in auto-vacuum or incremental vacuum mode
            # user version
            # incremental vacuum mode
            # application id
            # version-valid-for
            # sqlite version number
        ])

unittest.main()
