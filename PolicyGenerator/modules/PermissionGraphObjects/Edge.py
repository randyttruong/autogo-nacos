# Edge.py (DEPRECATED) 
# 
# Helper class for creating an `Edge` object 
# 
# Randy Truong 
# Northwestern University 
# 11 February 2024 

# Python library imports 
from enum import Enum 
from typing import defaultdict, Optional

# Local imports
from modules.PermissionGraphObjects import Node

class RequestType(Enum): 
    SERVICE_REQUEST    = 1 
    DATABASE_REQUEST   = 2

class EdgeType(Enum): 
    BELONGING_EDGE   = 1 
    REQUEST_EDGE     = 2 

class Edge: 
    def __init__(self, src: Node, dst: Node, edgeType: int, requestType: int) -> None:  
        self.src = src 
        self.dst = dst 
        self.edgeType = edgeType 
        self.requestType = requestType
        return None 

    def __eq__(self, other) -> bool: 
        if not isinstance(other, Edge): return False 

        return (self.src == other.src
                and self.dst == other.dst 
                and self.edgeType == other.edgeType
                and self.requestType == other.requestType)

# Edge.parseRequestName()
def parseRequestName(r: defaultdict) -> str: 
    pass

# Edge.parseRequestType()
def parseRequestType(rName: str) -> int: 
    match ("db" in rName.lower()): 
        case True: 
            return RequestType.SERVICE_REQUEST

        case False: 
            return RequestType.DATABASE_REQUEST

print(RequestType.SERVICE_REQUEST == 1)



