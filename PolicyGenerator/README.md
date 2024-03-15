# AutoArmor Policy Generator 
This is a program that automatically models microservice
architectures and automatically generates Kubernetes
microsegmentation policies in Python. 

## Dependencies
Python >= 3.11 
pytest == 8.0.2
PyYAML == 6.0.1

## File Descriptions 
`main.py` - Primary file to use Policy Generator Engine 
`setup.py` - Setup file for enabling absolute pathing + enabling the 

## Installation 
1. Install all dependencies
```
$ pip install -r requirements.txt
```

2. Run `setup.py` to enable local modules 
```
$ python3 setup.py install --user 
```

## Usage 
1. Run `main.py` with manifest paths 
```
$ python3 main.py <json_1> <json_2> ... 
```

Here is an example: 
- Unix-Based
```
$ python3 main.py ./example-json/callerService.json
```

- Windows 
```
$ python3 main.py -w .\example-json\callerService.json
```

## Inputs 
JSON files 

## Outputs 
YAML files 
