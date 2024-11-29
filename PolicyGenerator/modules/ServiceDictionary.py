# ManifestParser.py: AutoArmor JSON manifest file for Automated 
# Policy Generation for "Automatic Policy Generation for Inter-Service 
# Access Control of Microservices" 
# 
# Description
# This is a file for parsing a JSON-based manifest file. It simply just 
# returns a defaultdict object that represents the JSON-based manifest 
# file.
# 
# Randy Truong 
# Northwestern University 
# 1 November 2024 

# Python library import 
import os 
import sys 
import json 
import requests 
from collections import defaultdict
from typing import *  # Type hints for Python 
from enum import Enum 

class ServiceDictionary: 
    def __init__(self) -> None:  
        """
        params: 
        - 
        """
        self.services: defaultdict = {}
        self.baseUrl = "localhost:8080"  
        self.nacosQuery = "/nacos/v1/ns/instance/list"
        return None    
    
    def pingNacos(self): 
        """     
        A function that pings the Nacos instance registry
        params: 
        - 
        """    

        url = self.baseUrl + self.nacosQuery   
        headers = {} # TODO: Do I need any headers?  
        params = {}  # TODO: Do I need any parameters 

        resp = requests.get(url, params, headers, timeout=5) 
        if resp.status_code == 200:  
            print(f"Data: {resp.json()}")

            return resp.json() 
        else: 
            print(f"Error while attempting to access Nacos instance registry: {response.status_code}")  
