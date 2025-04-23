import json
import sys
import numpy as np
from scipy.spatial import ConvexHull


class NamedLocation:
    def __init__(self, id: str, lng: float, lat: float):
        self.id = id
        self.lng = lng
        self.lat = lat


def augment(filepath: str):
    with open(filepath) as fl:
        data = json.load(fl)

    places: list[NamedLocation] = []
    for section in data["sections"]:
        for place in section["places"]:
            places.append(NamedLocation(place["id"], lng=place["lng"], lat=place["lat"]))

    coords = np.array([(place.lng, place.lat) for place in places])
    hull = ConvexHull(coords)
    center = coords.mean(0)

    data["meta"]["center"] = {"lng": round(center[0], 5), "lat": round(center[1], 5)}

    bounding_indices = np.unique(hull.simplices.flat)

    data["meta"]["bounds"] = []
    bounding_places = [places[i] for i in bounding_indices]
    for _, place in enumerate(bounding_places):
        print(place.id, place.lat, place.lng)
        data["meta"]["bounds"].append({"id": place.id, "lat": place.lat, "lng": place.lng})

    # close the polygon
    data["meta"]["bounds"].append(data["meta"]["bounds"][0])

    return data


if __name__ == "__main__":
    try:
        file = sys.argv[1]
        augmented_json = augment(file)
        with open(file, "w") as f:
            json.dump(augmented_json, f, indent=4, ensure_ascii=False)
    except Exception as e:
        print(f"error: {e}")
        sys.exit(1)
