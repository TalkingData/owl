# coding:utf-8
from itertools import *

class Round_Robin():
    def __init__(self, data):
        self.data = data
        self.data_rr = self.get_item()

    def cycle(self, iterable):
        saved = []
        for element in iterable:
            yield element
            saved.append(element)

        while saved:
            for element in saved:
                yield element

    def get_item(self):
        count = 0
        for item in self.cycle(self.data):
		count += 1
		yield(count, item)

    def get_next(self):
        return self.data_rr.next()
