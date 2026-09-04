# Project Rules & Architecture: BVR AI & MCP Ecosystem

## Overview
This repository utilizes a modular, local-first architecture powered by **FastMCP**, **Antigravity**, **bvr-cli**, and **Gemini Spark**. The system bridges local terminal workflows with remote agentic execution, allowing seamless control across devices (including mobile-to-home workflows).

---

## 1. Core Architecture & Components
* **FastMCP Server (Provider):** The centralized engine defining custom tools and data interfaces. Built using Python/TypeScript via FastMCP network transports.
* **Hosting Nodes:**
  * *Local Development:* Ran locally during active feature development.
  * *Production / Home Lab:* Hosted on the idle Dell laptop running **CachyOS**, running containerized via Docker and exposed via secure public HTTPS tunnels (Cloudflare/ngrok).
* **Clients / Entry Points:**
  * **Gemini Spark (Web & Mobile):** Linked via Custom Connected Apps using the public MCP endpoint for on-the-road remote control.
  * **Antigravity:** Integrated locally via workspace/global `mcp_config.json` rules for autonomous agentic tool calling.
  * **bvr-cli:** Terminal-native client layer capable of programmatically invoking MCP commands.

---

## 2. MCP Configuration (`mcp_config.json`)
Antigravity and local clients reference the MCP server using standard configurations:
* **Remote/Tunnel Setup:**
  ```json
  {
    "mcpServers": {
      "cachy-mcp": {
        "serverUrl": "https://<your-tunnel-url>/mcp"
      }
    }
  }