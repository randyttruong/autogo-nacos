# parseTest.py
# This is a file for testing the 
# json parser file. 
#
# Randy Truong
# Northwestern University
# 2 March 2024 

import pytest 
from modules import ManifestParser
import unittest


parent = "../example-json/"
json1 = parent + "ratings.json"
json2 = parent + "reviews.json"
json3 = parent + "productpage.json"

dict1 = {
          "service": "ratings",
          "requests": [
            {
              "port": "3306",
              "name": "mysqldb",
              "type": "tcp",
              "url": "tcp://mysqldb:3306"
            },
            {
              "port": "27017",
              "name": "mongodb",
              "type": "tcp",
              "url": "mongodb://mongodb:27017/test"
            }
          ],
          "version": "v2"
        }

dict2 = {
  "service": "reviews",
  "requests": [
    {
      "path": "/ratings/*",
      "method": "GET",
      "name": "ratings",
      "type": "http",
      "url": "http://ratings:9080/ratings/*"
    }
  ],
  "version": "v1"
}

dict3 = {
          "service": "productpage",
          "requests": [
            {
              "path": "/ratings/*",
              "method": "GET",
              "name": "ratings",
              "type": "http",
              "url": "http://ratings:9080/ratings/*"
            },
            {
              "path": "/details/*",
              "method": "GET",
              "name": "details",
              "type": "http",
              "url": "http://details:9080/details/*"
            },
            {
              "path": "/reviews/*",
              "method": "GET",
              "name": "reviews",
              "type": "http",
              "url": "http://reviews:9080/reviews/*"
            }
          ],
          "version": "v1"
        }

class file1(unittest.TestCase): 
    def test1(self):
        parser = ManifestParser.ManifestParser(json1)
        parser.parse()

        self.assertEqual(parser.finalDict, dict1)
        self.assertEqual(parser.finalDictList, [dict1])

    def test2(self): 
        parser = ManifestParser.ManifestParser(json2)
        parser.parse()

        self.assertEqual(parser.finalDict, dict2)

    def test3(self): 
        parser = ManifestParser.ManifestParser(json3)
        parser.parse()

        self.assertEqual(parser.finalDict, dict3)

    def test4(self):
        self.maxDiff = None
        parser = ManifestParser.ManifestParser(json1)
        parser.parse()

        parser.setManifest(json2)
        parser.parse()

        sampleList = [dict1, dict2]
        assert sampleList == parser.finalDictList

        self.assertEqual(parser.finalDictList, [dict1, dict2])

    def test5(self):
        self.maxDiff = None
        parser = ManifestParser.ManifestParser(json1)
        parser.parse()

        parser.setManifest(json2)
        parser.parse()

        parser.setManifest(json3)
        parser.parse()

        sampleList = [dict1, dict2, dict3]
        assert sampleList == parser.finalDictList

        self.assertEqual(parser.finalDictList, [dict1, dict2, dict3])

    def test6(self): 
        parser = ManifestParser.ManifestParser(json3)
        parser.parse()

        parser.setManifest(json2)
        parser.parse()

        parser.setManifest(json1)
        parser.parse()

        sampleList = [dict3, dict2, dict1]

        assert sampleList == parser.finalDictList 

        self.assertEqual(parser.finalDictList, [dict3, dict2, dict1])



    def test_isupper(self):
        self.assertTrue('FOO'.isupper())
        self.assertFalse('Foo'.isupper())

    def test_split(self):
        s = 'hello world'
        self.assertEqual(s.split(), ['hello', 'world'])
        # check that s.split fails when the separator is not a string
        with self.assertRaises(TypeError):
            s.split(2)



