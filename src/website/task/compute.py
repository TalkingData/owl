# coding:utf-8
from __future__ import division
import numpy as np

class Stats:
	def __init__(self, sequence):
		self.sequence = np.array([float(item) for item in sequence])

	def sum(self):
		if len(self.sequence) < 1:
			return None
		else:
			return np.sum(self.sequence)

	def count(self):
		return self.sequence.shape[0]

	def min(self):
		if len(self.sequence) < 1:
			return None
		else:
			return np.amin(self.sequence)

	def max(self):
		if len(self.sequence) < 1:
			return None
		else:
			return np.amax(self.sequence)

	def avg(self):
		if len(self.sequence) < 1:
			return None
		else:
			return np.average(self.sequence)
