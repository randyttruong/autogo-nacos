# PolicyGenTest.py 
# This is a file for testing the policy generation 
# engine. 

import pytest 
from modules import ManifestParser, PermissionGraphEngine, PolicyGeneratorEngine
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

class policyGen(unittest.TestCase): 
	def test1(self): 
		parser = ManifestParser.ManifestParser(json1)
		parser.parse()

		parser.setManifest(json2)
		parser.parse()

		g = PermissionGraphEngine.PermissionGraph(parser.finalDictList)

		for i in range(len(g.manifests)): 
			g.mapServiceNode(g.manifests[i]["service"])
			g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])


		# for i in g.serviceGraph: 
			# print(i)

		for i in range(len(g.manifests)): 
			g.mapRequests(g.manifests[i]["service"], g.manifests[i]["requests"])

		n1 = Node.Node("ratings", 1, 1)
		n2 = Node.Node("reviews", 1, 1)

		poli = PolicyGeneratorEngine.PolicyGenerator(g)

		poli.getEdges()

		# for e in poli.rEdges: 
			# print(e.src.serviceName, "->", e.dst.serviceName)

		print("\n")
		poli.genEgressDenyAll();

		poli.genEgress(poli.rEdges)


		self.assertEqual(0,0)
















