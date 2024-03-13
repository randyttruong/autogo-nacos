# PermissionGraphTest.py 
# This is a file for testing the permission 
# graph generator 

import pytest 
from modules import ManifestParser, PermissionGraphEngine
from modules.PermissionGraphObjects import Edge, Node
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

class graphGen1(unittest.TestCase): 

	def test1(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		g = PermissionGraphEngine.PermissionGraph([dict1])

		serviceName = g.manifests[0]["service"]

		for i in g.manifests: 
			g.mapServiceNode(i["service"])

		self.assertEqual(i["service"], serviceName)

		n1 = Node.Node(serviceName, 1, 1)

		self.assertEqual(g.serviceGraph["ratings"], 
			n1)

	def test1_5(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		serviceName = g.manifests[0]["service"]

		for i in g.manifests: 
			g.mapServiceNode(i["service"])

		n1 = Node.Node(serviceName, 1, 1)

		self.assertEqual(g.serviceGraph["ratings"], 
			n1)
		self.assertEqual(1, 1)

	def test2(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse() 

		parser.setManifest(json2)
		parser.parse()
		self.maxDiff = None

		parser.setManifest(json3)
		parser.parse()
		self.assertEqual(parser.finalDictList, [dict1, dict2, dict3])


		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		s1 = g.manifests[0]["service"] # everything is sorted alphabetically lol 
		s2 = g.manifests[1]["service"]
		s3 = g.manifests[2]["service"]
		for i in range(len(g.manifests)):
			if i == 0: 
				self.assertEqual(s1, g.manifests[0]["service"])
			if i == 1: 
				self.assertEqual(s2, g.manifests[1]["service"])
			g.mapServiceNode(g.manifests[i]["service"])


		n1 = Node.Node(s1, 1, 1)
		n2 = Node.Node(s2, 1, 1)
		n3 = Node.Node(s3, 1, 1)


		self.assertEqual(g.serviceGraph["productpage"].serviceName, s1)
		self.assertEqual(g.serviceGraph["ratings"].serviceName, s2)
		self.assertTrue(g.serviceGraph["reviews"].serviceName, s3)


	def test3(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse() 

		parser.setManifest(json2)
		parser.parse()
		self.maxDiff = None

		parser.setManifest(json3)
		parser.parse()
		self.assertEqual(parser.finalDictList, [dict1, dict2, dict3])


		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		s1 = g.manifests[0]["service"] # everything is sorted alphabetically lol 
		s2 = g.manifests[1]["service"]
		s3 = g.manifests[2]["service"]
		for i in range(len(g.manifests)):
			if i == 0: 
				self.assertEqual(s1, g.manifests[0]["service"])
			if i == 1: 
				self.assertEqual(s2, g.manifests[1]["service"])
			g.mapServiceNode(g.manifests[i]["service"])


		n1 = Node.Node(s1, 1, 1)
		n2 = Node.Node(s2, 1, 1)
		n3 = Node.Node(s3, 1, 1)


		self.assertEqual(g.serviceGraph["productpage"], n1)
		self.assertEqual(g.serviceGraph["ratings"], n2)
		self.assertTrue(g.serviceGraph["reviews"], n3)

	def test4(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		parser.setManifest(json2) 
		parser.parse() 

		self.assertEqual(parser.finalDictList, [dict1, dict2]) 

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		s1 = g.manifests[0]["service"] # productpage
		s2 = g.manifests[1]["service"] # ratings 

		for i in range(len(g.manifests)): 
			g.mapServiceNode(g.manifests[i]["service"])
			g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])


		n1 = Node.Node(s1, 1, 1)
		n2 = Node.Node(s2, 1, 1)

		n1_2 = Node.Node(s1, 2, 1) # create corresponding ver nodes 
		n2_2 = Node.Node(s2, 2, 1 ) # create corresponding ver nodes 

		e1 = Node.Edge(n1, n1_2, 1, -1)
		e2 = Node.Edge(n2, n2_2, 1, -1)

		bAdjList1 = {"v2" : e1}
		bAdjList2 = {"v1": e2}

		self.assertEqual(g.serviceGraph["ratings"].serviceName, n1.serviceName)
		self.assertEqual(g.serviceGraph["ratings"].bAdjList, bAdjList1)

		self.assertEqual(g.serviceGraph["reviews"].serviceName, n2.serviceName)
		self.assertEqual(g.serviceGraph["reviews"].bAdjList, bAdjList2)

	def test5(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		parser.setManifest(json2) 
		parser.parse() 

		parser.setManifest(json3)
		parser.parse()

		self.assertEqual(parser.finalDictList, [dict1, dict2, dict3]) 

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		s1 = g.manifests[0]["service"] # productpage
		s2 = g.manifests[1]["service"] # ratings 
		s3 = g.manifests[2]["service"] # reviews 

		for i in range(len(g.manifests)): 
			g.mapServiceNode(g.manifests[i]["service"])
			g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])

		n1 = Node.Node(s1, 1, 1)
		n2 = Node.Node(s2, 1, 1)
		n3 = Node.Node(s3, 1, 1)

		n1_2 = Node.Node(s1, 2, 1)
		n2_2 = Node.Node(s2, 2, 1)
		n3_2 = Node.Node(s3, 2, 1)

		e1 = Node.Edge(n1, n1_2, 1, -1)
		e2 = Node.Edge(n2, n2_2, 1, -1)
		e3 = Node.Edge(n3, n3_2, 1, -1)

		bAdjList1 = { "v1" : e1 }
		bAdjList2 = { "v2" : e2 }
		bAdjList3 = { "v1" : e3 }

		self.assertEqual(g.serviceGraph["productpage"].serviceName, n1.serviceName)
		self.assertEqual(g.serviceGraph["productpage"].bAdjList, bAdjList1)

		self.assertEqual(g.serviceGraph["ratings"].serviceName, n2.serviceName)
		self.assertEqual(g.serviceGraph["ratings"].bAdjList, bAdjList2)

		self.assertEqual(g.serviceGraph["reviews"].serviceName, n3.serviceName)
		self.assertEqual(g.serviceGraph["reviews"].bAdjList, bAdjList3)

	def test6(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		self.assertEqual(parser.finalDictList, [dict1]) 

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		s1 = g.manifests[0]["service"] # ratings

		for i in range(len(g.manifests)): 
			g.mapServiceNode(g.manifests[i]["service"])
			g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])

		print(g.manifests[0]["version"])

		n1 = Node.Node(s1, 1, 1)
		n2 = Node.Node(s1, 2, 1) # version node of node 1 

		e = Node.Edge(n1, n2, 1, -1)

		bAdjList = {"v2" : e}

		n1.bAdjList = bAdjList

		self.assertEqual(g.serviceGraph["ratings"].serviceName, n1.serviceName)
		self.assertEqual(g.serviceGraph["ratings"].bAdjList, bAdjList)


	def test7(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		parser.setManifest(json2)
		parser.parse()

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		for i in range(len(g.manifests)): 
			g.mapServiceNode(g.manifests[i]["service"])
			g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])

		for i in range(len(g.manifests)): 
			g.mapRequests(g.manifests[i]["service"], g.manifests[i]["requests"])


		n1 = Node.Node("ratings", 1, 1)
		n2 = Node.Node("reviews", 1, 1)

		e1 = Node.Edge(n2, n1, 2, 2)

		rAdjList1 = { } 
		rAdjList2 = { "ratings" : e1 }

		self.assertEqual(g.serviceGraph["ratings"].rAdjList, rAdjList1)
		assert g.serviceGraph["reviews"].rAdjList["ratings"].src.serviceName == rAdjList2["ratings"].src.serviceName
		assert g.serviceGraph["reviews"].rAdjList["ratings"].dst.serviceName == rAdjList2["ratings"].dst.serviceName


	def test8(self): 
		self.assertEqual(0,0)

	def test9(self): 
		self.assertEqual(0,0)

	def test10(self):
		self.assertEqual(0,0)





































