#!/usr/bin/env python

import requests
import json


response = requests.post(
    'http://localhost:8080/message',
    {
        'from': 'from@blah.com',
        'subject': 'testing',
        'text': 'blah blah',
        'to': ['to@blah.com'],
    }
)

print response
print json.dumps(response.json(), indent=1)