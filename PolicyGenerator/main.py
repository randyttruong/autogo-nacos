# Automated Microsegmentation Policy Generation 
# for Cloud-Native Applications Using Service Discovery
# 
# Randy Truong, Javier Cuadra, David Hu (+ NU LIST) 
# Northwestern University 
# 10 February 2024 

# Necessary Imports
import os 
import sys 
import json 
from collections import defaultdict 
from typing import * 

# Custom Imports :)
from modules.ManifestParser import ManifestParser
from modules.PermissionGraphEngine import PermissionGraph
from modules.PolicyGeneratorEngine import PolicyGenerator



# import xml.etree.ElementTree as ET 

def main(filenames: List[str]): 
    parser = ManifestParser("")


    for filename in filenames: 
        parser.setManifest(filename)
        parser.parse()

    g = PermissionGraph(parser.finalDictList)

    for i in range(len(g.manifests)): 
        g.mapServiceNode(g.manifests[i]["service"])
        g.mapVersionNode(g.manifests[i]["service"], g.manifests[i]["version"])

    for i in range(len(g.manifests)):
        g.mapRequests(g.manifests[i]["service"], g.manifests[i]["requests"])

    pg = PolicyGenerator(g)

    pg.getEdges()

    pg.genEgress(pg.rEdges)



    # print(parser.finalDict["requests"])

if (__name__ == "__main__"):  
    if (len(sys.argv) < 2): 
        raise Exception("[ERROR]: Please add a valid input file\nUsage: python3 main.py <filename>\n") 

    filenames = sys.argv[1:]
    print(filenames)

    # ... 

    main(filenames)



