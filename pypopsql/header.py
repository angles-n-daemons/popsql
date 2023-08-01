from dataclasses import dataclass
from enum import Enum
from util import b2i

class FileFormatVersion(Enum):
    LEGACY = 1
    WAL = 2

class SchemaFormat(Enum):
    FORMAT_1 = 1
    FORMAT_2 = 2
    FORMAT_3 = 3
    FORMAT_4 = 4

class TextEncoding(Enum):
    UTF_8 = 1
    UTF_16le = 2
    UTF_16be = 3

@dataclass
class Version:
    major: int
    minor: int
    patch: int

class Header:
    def __init__(
        self,
        data: bytes,
    ):
        header_str = data[:16]
        if header_str != b'SQLite format 3\x00':
            raise Exception('header string not found, result is ', header_str)
        
        self.page_size = b2i(data[16:18])

        self.file_format_write_version = FileFormatVersion(data[18])
        self.file_format_read_version = FileFormatVersion(data[19])

        self.page_end_reserved_space = data[20]
        
        self.maximum_embedded_payload_fraction = data[21]
        self.minimum_embedded_payload_fraction = data[22]
        self.leaf_payload_fraction = data[23]

        self.file_change_counter = b2i(data[24:28])

        self.db_size_in_pages = b2i(data[28:32])
        self.first_freelist_trunk_page = b2i(data[32:36])
        self.num_freelist_pages = b2i(data[36:40])

        self.schema_cookie = b2i(data[40:44])
        self.schema_format_number = SchemaFormat(b2i(data[44:48]))

        self.default_page_cache_size = b2i(data[48:52])
        self.largest_btree_root_page = b2i(data[52:56])
        self.text_encoding = TextEncoding(b2i(data[56:60]))
        self.user_version = b2i(data[60:64])

        self.incremental_vacuum_mode = b2i(data[64:68])
        self.application_id = b2i(data[68:72])
        self.version_valid_for = b2i(data[92:96])
        vb = b2i(data[96:100])
        major = vb // 1000000
        minor = (vb % 1000000) // 1000
        patch = vb % 1000
        self.version = Version(major, minor, patch)


    def _debug_print_values(self):
        print('page size', self.page_size)
        print('file format write version', self.file_format_write_version)
        print('file format read version', self.file_format_read_version)
        print('page reserved space', self.page_end_reserved_space)
        print('maximum embedded payload fraction', self.maximum_embedded_payload_fraction)
        print('minimum embedded payload fraction', self.minimum_embedded_payload_fraction)
        print('leaf payload fraction', self.leaf_payload_fraction)
        print('file change counter', self.file_change_counter)
        print('database size in pages', self.db_size_in_pages)
        print('first freelist trunk page', self.first_freelist_trunk_page)
        print('number of freelist pages', self.num_freelist_pages)
        print('schema cookie', self.schema_cookie)
        print('schema format number', self.schema_format_number)
        print('default page cache size', self.default_page_cache_size)

        print('page of largest btree root', self.largest_btree_root_page)
        print('user version', self.user_version)
        print('incremental_vacuum_mode', self.incremental_vacuum_mode)

        print('application id', self.application_id)
        print('version valid for', self.version_valid_for)
        print('sqlite version number', self.version)
