#!/usr/bin/env python3
"""Render the VyOS post-config bootup script from Jinja2 templates.

Required env vars: WAN_MAC, LAN_MAC, SSH_AUTHORIZED_KEY
Optional env var:  ENABLE_MDNS ("true" or "false", defaults to "false")

Usage:
    pip install jinja2
    WAN_MAC=... LAN_MAC=... SSH_AUTHORIZED_KEY=... python render_config.py
"""

import os
import sys
from dataclasses import asdict, dataclass
from pathlib import Path

from jinja2 import Environment, FileSystemLoader, StrictUndefined


@dataclass
class MDNSConfig:
    enabled: bool = False
    interface_mdns: str = ""  # ethX
    hostname: str = ""


@dataclass
class BaseConfig:
    wan_mac: str
    lan_mac: str
    ssh_authorized_key: str


def render_template(template_path: str, dest_path: str, context: dict) -> None:
    """Render the Jinja2 template at template_path into dest_path."""
    src = Path(template_path)
    dest = Path(dest_path)

    env = Environment(
        loader=FileSystemLoader(src.parent),
        undefined=StrictUndefined,  # fail loudly on missing variables
        keep_trailing_newline=True,
        trim_blocks=True,
        lstrip_blocks=True,
    )
    template = env.get_template(src.name)
    rendered = template.render(**context)

    dest.parent.mkdir(parents=True, exist_ok=True)
    dest.write_text(rendered)


def create_base_config_file(template_path: str, dest_path: str, config: BaseConfig) -> None:
    render_template(template_path, dest_path, asdict(config))


def create_mdns_config_file(template_path: str, dest_path: str, config: MDNSConfig) -> None:
    if not config.enabled:
        return
    # Exposed to the template as `mdns.enabled`, `mdns.interface_mdns`, `mdns.hostname`
    render_template(template_path, dest_path, {"mdns": asdict(config)})


def parse_env_var_bool(value: str) -> bool:
    if value not in ("true", "false"):
        raise ValueError("Invalid value. Only 'true' or 'false' are available")
    return value == "true"


def main() -> None:
    wan_mac = os.environ.get("WAN_MAC", "")
    lan_mac = os.environ.get("LAN_MAC", "")
    ssh_authorized_key = os.environ.get("SSH_AUTHORIZED_KEY", "")
    mdns_enable = os.environ.get("ENABLE_MDNS", "false")

    if not (wan_mac and lan_mac and ssh_authorized_key):
        sys.exit(
            "Environment variable not set. Make all three have a value: "
            "WAN_MAC, LAN_MAC, SSH_AUTHORIZED_KEY"
        )

    create_base_config_file(
        "templates/vyos-postconfig-bootup.script.tmpl",
        "scripts/vyos-postconfig-bootup.script",
        BaseConfig(
            wan_mac=wan_mac,
            lan_mac=lan_mac,
            ssh_authorized_key=ssh_authorized_key,
        ),
    )

    create_mdns_config_file(
        "scripts/vyos-postconfig-bootup.script",
        "scripts/vyos-postconfig-bootup.script",
        MDNSConfig(
            enabled=parse_env_var_bool(mdns_enable),
            interface_mdns="",
            hostname="",
        ),
    )


if __name__ == "__main__":
    main()
