from typing import BinaryIO

class Pager:
    def __init__(
        self,
        file_name: str,
        page_size: int = 4096,
    ):
        self.file = open(file_name, 'rb+')
        self.page_size = page_size

    def get_page(
        self,
        page_number: int,
    ) -> bytes:
        self.file.seek(self.get_offset(page_number))
        return self.file.read(self.page_size)

    def write_page(
        self,
        page_number: int,
        data: bytes,
    ):
        raise ValueError('not implemented')

    def get_offset(
        self,
        page_number: int,
    ):
        return (page_number - 1) * self.page_size
