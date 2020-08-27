#!/usr/bin/env python

import requests
import json


response = requests.put(
    "http://localhost:8001/event",
    {
        "type": "Planned Event",
        "contact": "gmendiola@connexity.com",
        "date": "2020-08-19Z22:39:06",
        "summary": "test summary",
        "body": "test body",
        "impact": "test impact",
        "other": "test other",
    },
)

print(response)
print(json.dumps(response.json(), indent=1))
