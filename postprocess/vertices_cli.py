import json
import sys

import numpy as np
from scipy.spatial import ConvexHull

try:
    file = sys.argv[1]
except:
    print("No file specified")
    sys.exit(1)

with open(file) as f:
    data = json.load(f)

places: list[tuple[str, float, float]] = []
for section in data["sections"]:
    for place in section["places"]:
        places.append((place["id"], place["lng"], place["lat"]))

coords = np.array([(lng, lat) for _, lng, lat in places])
hull = ConvexHull(coords)
center = coords.mean(0)

print(f"center: {round(center[0], 5)}, {round(center[1], 5)}")

bounding_indices = np.unique(hull.simplices.flat)

bounding_places = [places[i] for i in bounding_indices]
for i, place in enumerate(bounding_places):
    print(f"vertex {i}, {place[0]}, lat: {place[2]}, lng: {place[1]}")
