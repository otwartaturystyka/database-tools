// Script to query Notion database and generate GeoJSON of locations
// Requires Node.js and the @notionhq/client package
// Installation: npm install -g @notionhq/client

import { Client } from "@notionhq/client";

// Initialize Notion client
// Replace with your own integration token
const notion = new Client({
  auth: process.env.NOTION_INTEGRATION_TOKEN,
});

const databaseId = "3fe6c207c85f4ccfb4ae5f9649848899";

// Query the database
const response = await notion.databases.query({
  database_id: databaseId,
  filter: {
    property: "Koordynaty",
    rich_text: {
      is_not_empty: true,
    },
  },
});

// Transform the Notion data into GeoJSON format
const features = response.results
  .map((page) => {
    const nazwa =
      page.properties.Nazwa?.title[0]?.plain_text || "Unnamed Location";
    const miejscowosc =
      page.properties.Miejscowość?.multi_select
        .map((item) => item.name)
        .join(", ") || "";
    const sekcja = page.properties.Sekcja?.select?.name || "";
    const status = page.properties.Status?.select?.name || "";

    // Parse coordinates
    const coordsText =
      page.properties.Koordynaty?.rich_text[0]?.plain_text || "";
    let lat = 0;
    let lng = 0;

    if (coordsText) {
      const coordParts = coordsText
        .split(",")
        .map((part) => parseFloat(part.trim()));
      if (coordParts.length === 2) {
        lat = coordParts[0];
        lng = coordParts[1];
      }
    }

    // Only include points with valid coordinates
    if (lat && lng) {
      let description = sekcja;
      if (miejscowosc) {
        description += ", " + miejscowosc;
      }
      if (status === "Już nie istnieje") {
        description += " (Już nie istnieje)";
      }

      return {
        type: "Feature",
        properties: {
          name: nazwa,
          description: description,
        },
        geometry: {
          type: "Point",
          coordinates: [lng, lat], // GeoJSON uses [longitude, latitude] order
        },
      };
    }
    return null;
  })
  .filter((feature) => feature !== null);

// Create the GeoJSON object
const geojson = {
  type: "FeatureCollection",
  features: features,
};

// Print the GeoJSON to stdout
console.log(JSON.stringify(geojson, null, 2));
