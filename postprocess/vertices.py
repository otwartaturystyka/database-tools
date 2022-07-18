import json
import sys
import numpy as np
from scipy.spatial import ConvexHull


def augument(filepath: str):
    with open(filepath) as f:
        data = json.load(f)

    places: list[tuple[str, float, float]] = []
    for section in data["sections"]:
        for place in section["places"]:
            places.append((place["id"], place["lng"], place["lat"]))

    coords = np.array([(lng, lat) for _, lng, lat in places])
    hull = ConvexHull(coords)
    center = coords.mean(0)

    data["center"] = {"lng": round(center[0], 5), "lat": round(center[1], 5)}

    bounding_indices = np.unique(hull.simplices.flat)

    data["bounds"] = []
    bounding_places = [places[i] for i in bounding_indices]
    for _, place in enumerate(bounding_places):
        lat = round(place[1], 5)
        lng = round(place[2], 5)
        data["bounds"].append({"id": place[0], "lng": lng, "lat": lat})

    # close the polygon
    data["bounds"].append(data["bounds"][0])

    return data


if __name__ == "__main__":
    try:
        file = sys.argv[1]
        augumented_json = augument(file)
        with open(file, "w") as f:
            json.dump(augumented_json, f, indent=4, ensure_ascii=False)
    except:
        print("No file specified")
        sys.exit(1)
