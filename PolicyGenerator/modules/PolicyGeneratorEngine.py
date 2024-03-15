# PolicyGeneratorEngine.py: AutoArmor Policy Generation Phase 
#
# Automated Microsegmentation Policy Generation 
# for Cloud-Native Applications Using Service Discovery
#
# Given a directed acyclic graph (aka the Permission Graph) of a 
# a microservice application, generate microsegmentation policies 
# for increasing application defense. Note that the relevant 
# policies are all going to be ingress policies 
# 
# 
# Randy Truong
# Northwestern University 
# 10 February 2024 



# Here is an example of the auto-generated .yaml files 
# apiVersion: networking.k8s.io/v1
# kind: NetworkPolicy
# metadata:
#   name: quickstart-fourth
# spec:
#   egress:
#   - to:
#     - podSelector:
#         matchLabels:
#          app.consumer: consumer
#     - podSelector:
#         matchLabels:
#          app: third
#  podSelector:
#    matchLabels:
#      app: fourth
#  policyTypes:
#  - Egress

# Python library imports 
import os 
import sys 
import json
import yaml
from collections import defaultdict 
from typing import * 

# Local imports 
from modules.PermissionGraphEngine import PermissionGraph
from modules.PermissionGraphObjects.Node import Node, Edge, parseRequestName, parseRequestType, RequestType

class PolicyGenerator: 
    def __init__(self, g: PermissionGraph) -> None: 
        """
        params: 
        - serviceCalls: An array of string arrays 

        """

        self.finalPolicyTemplate: defaultdict = {  
                                         'apiVersion' : 'networking.k8s.io/v2', 
                                         'policyType' : 'NetworkPolicy', 
                                         'metadata'   : {"name" : ""} ,
                                         } 

        self.finalPolicy = self.finalPolicyTemplate

        self.ingress: defaultdict = { "ingress" : None }                   
        self.egress: defaultdict = { "egress" : None }    

        self.graph: PermissionGraph = g 
        self.bEdges: List[Edge] = []
        self.rEdges: List[Edge] = []
        return None 

    # PolicyGenerator.generateName()

    # PolicyGenerator.getEdges() 
    def getEdges(self) -> None: 
        """ 
        params: 
        - self.graph: 
        - self.bEdges: 
        - self.rEdges: 
        """ 

        # For service in the serviceGraph 
        for (service, value) in self.graph.serviceGraph.items(): 
            bAdjList: defaultdict = value.bAdjList # Grab Belonging Edges 
            rAdjList: defaultdict = value.rAdjList # Grab Request Edges 

            # Add the belonging edge to the self.bEdges list 
            for (targetService, edge) in bAdjList.items(): 
                self.bEdges.append(edge)

            # Add the request edge to the self.rEdges list 
            for (targetService, edge) in rAdjList.items(): 
                self.rEdges.append(edge)

        pass 


    # PolicyGenerator.genIngress() - (DEPRECATED (?))
    def genIngress(self, e: Edge) -> None: 
        """ 
        params: 

        desc: 
        - Generates an ingress policy for the  

        Note: There are four different types of selectors 
        for policy generation: 
            1. podSelector 
            2. 
            3. 
            4. 



        ingress:
            - from:
            - namespaceSelector:
        matchLabels:
            user: alice
        podSelector:
            matchLabels:
                role: client
        """ 
        src = e.src
        dst = e.dst



        # ingress = { "ingress" : None } 
        # fromBlock = { "from" : None }
        # ipBlock = { "ipBlock" : None }
        # nsBlock = { "- namespaceSelector" : None } 
        pass 

    # PolicyGenerator.genIngressDenyAll()
    def genIngressDenyAll(self) -> None: 
        """ 
        params: 
        desc: 

        apiVersion: networking.k8s.io/v1
        kind: NetworkPolicy
        metadata:
            name: default-deny-ingress
        spec:   
            podSelector: {}
            policyTypes:
              - Ingress        
        """ 

        # Reset the policy template 
        self.finalPolicy = self.finalPolicyTemplate

        # Init values 
        self.finalPolicy["metadata"]["name"] = "default-deny-ingress"
        spec = { "podSelector" : {} }
        spec["policyTypes"] = ["Ingress"]
        self.finalPolicy["spec"] = spec 
        pass 

    # PolicyGenerator.genIngressAcceptAll()
    def genIngressAcceptAll(self) -> None: 
        """
        params: 
        desc: 

        apiVersion: networking.k8s.io/v1
        kind: NetworkPolicy
        metadata:
            name: allow-all-ingress
        spec:
            podSelector: {}
            ingress:
            - {}
            policyTypes:
            - Ingress
        """ 

        # Reset the policy tempalte 
        self.finalPolicy = self.finalPolicyTemplate

        # Init values 
        self.finalPolicy["metadata"]["name"] = "allow-all-ingress"
        spec = { "podSelector" : {} } 
        spec["ingress"] = {} 
        spec["policyTypes"] = { "ingress": None} 
        self.finalPolicy["spec"] = spec 
        pass 

    # TODO PolicyGenerator.genPodSelector()
    def genPodSelector(self, e: Edge) -> None: 
        pass 

    # TODO PolicyGenerator.genEgress()
    def genEgress(self, E: List[Edge], windows: bool = False) -> None: 
        """
        params: 

        desc: 
        - Generates an egress policy. 

        Note: There are four different types of selectors: 
            1. podSelector selects 
            2. namespaceSelector selects 
            3. namespaceSelector + podSelector selects 
            4. ipBlock selects 

        Note that there are two different types of matchLabels: 
        1. If the service is a consumer/provider, 
        then use app.{consumer/provider}: {consumer/provider}
        2. Otherwise, just use app: {version of service}

        TODO 
        - Add podSelector entry for app.consumer 
        - Add podSelector entry for app.provider 
        - Add podSelector entries for non-consumers/non-providers 
        """

        # Grab `src` and `dst` from Edge `e` 
        for e in E: 
            src: Node = e.src
            dst: Node = e.dst 

            self.finalPolicy = self.finalPolicyTemplate

            self.finalPolicy["metadata"]["name"] = f"{src.serviceName}-microseg-policy"

            spec = { "spec" : None } 
            egress = [ 
                    { "to" : 
                    [{ "podSelector" : {"matchLabels" : 
                    {"serviceName" : dst.serviceName}}}]}
                ]

            to = [] # TODO, start working from to 

            podSelector = { "podSelector" : { "matchLabels" : 
                                        { "serviceName": src.serviceName }}} 
            matchlabels = { "matchLabels" : None } 


            spec["spec"] = podSelector
            spec["spec"]["egress"] = egress 
            spec["spec"]["policyTypes"] = ["Egress"]


            self.finalPolicy["spec"] = spec["spec"]

            if not windows: 
                self.dump(f"./outputs/{src.serviceName}_{dst.serviceName}.yaml", self.finalPolicy)
            else: 
                self.dump(f".\\outputs\\{src.serviceName}_{dst.serviceName}.yaml", self.finalPolicy)
            pass

    # PolicyGenerator.genEgressDenyAll()
    def genEgressDenyAll(self) -> None: 
        """ 
        params: 
        desc: 

        apiVersion: networking.k8s.io/v1
            kind: NetworkPolicy
        metadata:
            name: default-deny-egress
        spec:
            podSelector: {}
            policyTypes:
            - Egress    
        """ 

        # Reset the policy template 
        self.finalPolicy = self.finalPolicyTemplate

        # Init values 
        self.finalPolicy["metadata"]["name"] = "default-deny-egress"
        spec = { "podSelector" : {} }
        spec["policyTypes"] = ["Egress"]
        self.finalPolicy["spec"] = spec 
        pass 

    # PolicyGenerator.genEgressAcceptAll()
    def genEgressAcceptAll(self) -> None: 
        """
        params: 
        desc: 

        apiVersion: networking.k8s.io/v1
        kind: NetworkPolicy
        metadata:
            name: allow-all-egress
        spec:
            podSelector: {}
            egress:
            - {}
            policyTypes:
            - Egress
        """ 

        # Reset the policy tempalte 
        self.finalPolicy = self.finalPolicyTemplate

        # Init values 
        self.finalPolicy["metadata"]["name"] = "allow-all-ingress"
        spec = { "podSelector" : {} } 
        spec["ingress"] = {} 
        spec["policyTypes"] = { "ingress": None} 
        self.finalPolicy["spec"] = spec 
        pass 


    # PolicyGenerator.genDenyAllIngressEgress()
    def genDenyAllIngressEgress(self) -> None: 
        """
        params: 
        desc: 

        apiVersion: networking.k8s.io/v1
        kind: NetworkPolicy
        metadata:
            name: default-deny-all
        spec:
            podSelector: {}
            policyTypes:
            - Ingress
            - Egress
        """

        # Reset policy template 
        self.finalPolicy = self.finalPolicyTemplate

        # Init values 
        spec = {} 
        spec["podSelector"] = {} 
        spec["policyTypes"] = { "Ingress" : None, 
                                "Egress": None } 

        spec.finalPolicy["spec"] = spec     
        pass 

    # TODO PolicyGenerator.generatePolicy()
    def generatePolicy(self) -> None: 
        """
        params: 
        - self.rEdges: 

        desc: 
        - Given the rEdges 
        """
        # Base Case: Not exists any serviceCalls 
        if (len(self.rEdges) == 0): 
            raise Exception("[ERROR]: No parsable service calls.")

        # For edge in request edges: 
        for e in self.rEdges: 

            # Generate an egress policy     
            self.genEgress(e)
            # self.genIngress(e)
            pass 

        pass 

    # TODO PolicyGenerator.dump()
    def dump(self, filename: str, yamlDict: defaultdict) -> None:
        """
        params: None

        desc: 
        - 
        """
        try:  
            data = yaml.dump(yamlDict, sort_keys=False)
            print(data)
            with open(filename, "a") as f: 
                f.write(data)
        except FileExistsError: 
            with open(filename, "w", encoding="utf-8") as f: 
                f.truncate() 
                yaml.dump(data=self.finalPolicy, stream=f, allow_unicode=True) 

        pass 


