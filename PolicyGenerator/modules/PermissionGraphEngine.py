# PermissionEngine.py: AutoArmor Permission Graph Generation Engine 
# 
# Automated Microservice Policy Generation 
# for Cloud-Native Applications Using Service 
# Discovery 
# 
# Given a parsed manifest file. Generate a directed acyclic 
# graph representation of a microservice architecture-
# based application. 
#
# Randy Truong 
# Northwestern University 
# 10 February 2024 

""" 
Example of rawGraph: 
    [ node1: { 
              "type":, 
              "version": , 
              "adjList": , 
              }, 
      node2: { 
              "type": , 
              "version": , 
              "adjList" , 
              }, 
      node3: {
              "type": , 
              "version": , 
              "adjList": , 
              }, 
      node4: { 
              "type": , 
              "version": , 
              "adjList" , 
              }, 
      node5: { 
              "type": , 
              "version: , 
              "adjList" , 
              }, 
    ] 

Example of serviceGraph:   
    [  
      "serviceNode1": {
                        "adjList": [] 
                       }, 
      "serviceNode2": {
                        "adjList": [] 
                       }, 
      "serviceNode3": { 
                        "adjList": [] 
                       }, 
      "serviceNode4": {
                        "adjList": [] 
                       }
     ] 
     
""" 

# Python library imports
import os 
import sys 
import json 
from collections import defaultdict
from typing import * 

# Custom imports 
from modules.PermissionGraphObjects.Node import Node, Edge, parseRequestName, parseRequestType, RequestType
from modules import ManifestParser

class PermissionGraph: 
    def __init__(self, manifests: List[defaultdict]) -> None: 
        """
        params:  
        - defaultdict manifestDict: This is a dictionary version 
        of the parsed manifest file. 
        
        returns: 
        - None 

        desc: 
        - Constructor for `PermissionGraph` object 

        attributes: 
        - `List[defaultdict] PermissionGraph.manifests`: 
                    List of json-ified manifest files that will
                    be turned into a permission graph. 

        - `defaultdict self.rawGraph`: 
                    List of all nodes N_{s} and N_{v} within G.

        - `defaultdict self.serviceGraph`: 
                    Hierarchal mapping of N_{s} and N_{v} within G.
        """
        manifests.sort(key = lambda a: a["service"]) 

        self.manifests: List[defaultdict] = manifests 
        # self.rawGraph: defaultdict = {} 
        self.serviceGraph: defaultdict = {} 
        self.sortedGraph = []
        return None 

    # PermissionGraph.mapServiceNode() 
    def mapServiceNode(self, service: str) -> None: 
        """
        params: 
        - None 

        desc:
        Creates a service node as well as creates a corresponding 
        """
        # Init Service Node Object 
        currentNode = Node(service, 1) # 1 == SERVICE_NODE
        self.serviceGraph[service] = currentNode
        pass

    # PermissionGraph.mapVersionNode()
    def mapVersionNode(self, service: str, version: str) -> None: 
        """

        desc: 
        Creates a version node for a given service node 
        """
        # Find corresponding service node
        try: 
            serviceNode = self.serviceGraph[service]
        except: 
            raise Exception(f"[DEBUG] Tried to add a version node\
                    without a service node {service}")


        # Init Version Node Object 
        versionNode = Node(service, 2, version) # 2 == VERSION_NODE
        self.serviceGraph[service].addBelongingEdge(versionNode)
        pass

    # TODO PermissionGraph.mapRequests() 
    def mapRequests(self, serviceName: str,  
                    requests: List[defaultdict]) -> None: 
        """
        params: 
        - str self.serviceName: 
                The name of the service that we are adding nodes to 
        - str self.version: 
                The version of the service that we are adding nodes to 
        - str self.requests: 
                The version of the service that we are adding nodes to 

        returns: 
        - None 

        desc: 
        This is a function that maps requests as one of two edges 
        in the permission graph: 
        - E_{b}: Belonging Edge (Edge that connects a version node to 
        a service node) 
        - E_{r}: Request Edge (Edge that connects a service to 
        another service) 
        """

        currentService = serviceName 

        try:
            currentService: Node = self.serviceGraph[serviceName]
        except: 
            raise ValueError(f"[DEBUG] No entry exists for {serviceName}")


        # TODO work on the request engine!!!
        for r in requests: 
            rName: str = r["name"]
            rType: int = parseRequestType(rName)

            # If not inter-service request -> continue 
            # 
            # TODO Will add database fnality later, since unsure how to 
            # add databases as nodes 
            if (rType == RequestType(1)): 
                continue

            # Otherwise, add an edge 
            try: 
                dst = self.serviceGraph[rName] # I should probably make an entry 
                                               # for things that don't actually exist 
                                               # within the graph 
                newEdge = currentService.addRequestEdge(dst, RequestType(2))

            except: 
                print(f"[DEBUG] Unable to find destination service for {rName}. Creating entry now...")

                newNode = Node(rName, 1)
                self.serviceGraph[rName] = newNode 
                newEdge = currentService.addRequestEdge(newNode, RequestType(2))


        pass 

    # PermissionGraph.generateNodes() 
    def generateNodes(self) -> None: 
        """ 
        params: 
        - None 

        returns: 
        - None 

        desc: 
        - Generates permission graph from `PermissionGraph.manifest` List[defaultdict]
        """ 

        # Remark: The permission graph is generated utilizing 
        # the following attributes: 
        # 
        # G = (N_{s}, N_{v}, E_{b}, E_{r}) where 
        # 
        #   N_{s} (Service Node):    
        #     Generic service for microservice architecture  
        #   N_{v} (Version Node):    
        #     Variation of N_{s} that may be invoked by other permissions 
        #   E_{b} (Belonging Edge):  
        #     An edge that connects N_{v} with its corresponding N_{s}
        #   E_{r} (Request Edge):    
        #     An edge that connects N_{s1} with another N_{s2} 
        # 
        # Plan: 
        #   1. Generate overarching service node
        #   2. Generate corresponding version nodes u

        manifests = self.manifests
        manifestLength = len(manifests)

        match manifestLength:

            # Case 1: len(manifests) == 0 -> Throw exception 
            case 0: 
                raise Exception("No manifests/services found.")

            # Case 2: len(manifests) == 1 -> Create service node and attach version + edges
            case 1: 
                service = self.manifests[0]
                self.mapServiceNode(service, True)
                self.mapVersionNode(service, False) 

            # Case 3: Create multiple service nodes 
            case _: 
                for m in manifests: 

                    service: str = m["service"] 
                    version: str = m["version"]
                    # Not in service graph -> Create service node + version node 
                    if (service not in self.serviceGraph): 
                        self.mapServiceNode(service) 
                        self.mapVersionNode(service, version) 

                    # Otherwise, create new version node 
                    else: 
                        self.mapManifest(service, False) 

    # PermissionGraph.generateEdges()
    def generateEdges(self) -> None: 
        """
        params: 
        - None 

        returns: 
        - None 

        desc: 
        - Generates edges for the permission graph from  
        `PermissionGraph.manifest` List[defaultdict]
        """

        manifests = self.manifests
        for m in manifests: 
            service: str = m["service"]
            version: str = m["version"]
            requests: str = m["requests"]

            self.mapRequests(service, version, requests)


    # PermissionGraph.dfs()
    def dfs(self, curr: defaultdict, visited: set, 
            path: set, finalSort: set) -> bool: 
        # Base Case 0: If visited -> return None 
        if (curr in visited): 
            return True 

        # Base Case 1: If in same path -> return False 
        if (curr in path): 
            return False 

        # Recurse through the rest of the graph 
        # for nei in adjList[curr]: 
        #     dfs(curr=, visited=, path= )

        
        return True 

        

    # TODO PermissionGraph.topSort()
    def topSort(self) -> None: 
        """ 
        params: 
        - None 

        returns: 
        - None 
        
        desc: 
        - Performs a topological sorting of the PermissionGraph by 
        prioritizing service nodes first, then their versions 
        """ 
        
        finalSort: List[tuple] = [] 
        path: set = set()
        visited: set = set()

        for node in graph: 
            if (not dfs(node, visited, path, finalSort)): 
                return False 

        return True 

    # TODO renderGraph()
    def renderGraph(self) -> None: 
        pass 


