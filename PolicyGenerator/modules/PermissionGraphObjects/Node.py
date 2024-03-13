# Objects.py
#
# Helper class for creating graph objects
# - Node 
# - Edge 
# 
# Additionally defines NodeType, EdgeType, and RequestType 
#
# Randy Truong 
# Northwestern University 
# 11 February 2024 

# Python library imports 
from __future__ import annotations # Allow class method to take type of self

from collections import defaultdict
from enum import Enum # Add Enum support 

class NodeType(Enum): 
    SERVICE_NODE   = 1 
    VERSION_NODE   = 2

class EdgeType(Enum): 
    BELONGING_EDGE   = 1 
    REQUEST_EDGE     = 2 

class RequestType(Enum): 
    NON_REQUEST        = -1 
    DATABASE_REQUEST    = 1 
    SERVICE_REQUEST   = 2


class Node: 
    # Node.__init__()
    def __init__(self, serviceName: str, nodeType: int, 
                 version: int = 1) -> None:
        """
        desc: 
        - Initializes a `Node` object for the `PermissionGraph` object 
        """
        try: 
            NodeType(nodeType)
        except: 
            raise ValueError(f"[ERROR] Invalid nodeType {nodeType}")

        self.serviceName = serviceName 
        self.type: NodeType = nodeType 
        self.version = version
        self.bAdjList: defaultdict = {} # (serviceName, Edge())
        self.rAdjList: defaultdict = {} # (serviceName, Edge()) 
        pass 

    # Node.__eq__()
    def __eq__(self, other) -> bool: 
        if not isinstance(other, Node): return False 

        return (self.serviceName == other.serviceName
                and self.type == other.type 
                and self.version == other.version 
                and self.bAdjList == other.bAdjList
                and self.rAdjList == other.rAdjList)

    # Node.addBelongingEdge()
    def addBelongingEdge(self, dst: Node) -> None: 
        """
        desc: 
        - Adds a belonging edge to `dst` from current node 
        """

        # If `dst.nodeType` is not VERSION_NODE -> Error 
        if (dst.type != 2):
            print(f"This is the type {dst.type}")
            raise Exception(f"[DEBUG] Attempted to add a non-version node {dst.serviceName} \
                    to {self.serviceName}") 

        newEdge = Edge(src=self, dst=dst, 
                       edgeType=1, requestType=-1) 

        self.bAdjList[dst.version] = newEdge # (version, Edge) 
        pass 

    # Node.addRequestEdge()
    def addRequestEdge(self, dst: Node, requestType: RequestType) -> None: 
        """
        desc: 
        - Adds a request edge to `dst` from current node 
        """

        if (dst.type != 1): 
            raise Exception(f"[DEBUG] Attempted to add a request edge \
                        to non-service node {dst.serviceName}")

        newEdge = Edge(src=self, dst=dst, 
                       edgeType=2, requestType=requestType)

        
        self.rAdjList[dst.serviceName] = newEdge # (dst., Edge)
        pass

class Edge: 
    def __init__(self, src: Node, dst: Node, edgeType: EdgeType, 
                 requestType: RequestType) -> None:  
        self.src = src 
        self.dst = dst 
        self.edgeType: EdgeType = edgeType 
        self.requestType = requestType
        return None 

    def __eq__(self, other) -> bool: 
        if not isinstance(other, Edge): return False 

        return (self.src.serviceName == other.src.serviceName 
                and self.dst.serviceName == other.dst.serviceName 
                and self.edgeType == other.edgeType
                and self.requestType == other.requestType)


# TODO 
# Node.parseRequestName()
def parseRequestName(r: defaultdict) -> str: 
    pass

# Node.parseRequestType()
def parseRequestType(rName: str) -> int: 
    match ("db" in rName.lower()): 
        case True: 
            return RequestType(1)

        case False: 
            return RequestType(2) 







