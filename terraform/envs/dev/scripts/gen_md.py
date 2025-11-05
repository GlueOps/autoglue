#!/usr/bin/env python3
import argparse, json, sys
from textwrap import dedent

def type_to_str(t):
    if isinstance(t, list):
        return "[" + ", ".join(type_to_str(x) for x in t) + "]"
    if isinstance(t, dict):
        if "attribute_types" in t:
            return "object(" + ", ".join(f"{k}:{type_to_str(v)}" for k,v in t["attribute_types"].items()) + ")"
        return json.dumps(t)
    return str(t)

def flags(spec):
    f = []
    if spec.get("required"): f.append("required")
    if spec.get("optional"): f.append("optional")
    if spec.get("computed"): f.append("computed")
    if spec.get("sensitive"): f.append("sensitive")
    return ", ".join(f) if f else "-"

def md_table_row(name, spec):
    t = spec.get("type")
    t_str = type_to_str(t) if t is not None else "(block)"
    desc = (spec.get("description") or "").replace("\n", " ")
    return f"| `{name}` | {t_str} | {flags(spec)} | {desc} |"

def render_block(block):
    out = []
    attrs = block.get("attributes") or {}
    if attrs:
        out += ["| Name | Type | Flags | Description |",
                "|---|---|---|---|"]
        for k, v in attrs.items():
            out.append(md_table_row(k, v))
        out.append("")
    # nested blocks
    for bname, b in (block.get("block_types") or {}).items():
        out.append(f"**Nested block `{bname}`** (mode: {b.get('nesting_mode','')})")
        nb = b.get("block", {})
        nattrs = nb.get("attributes") or {}
        if nattrs:
            out += ["| Name | Type | Flags | Description |",
                    "|---|---|---|---|"]
            for k, v in nattrs.items():
                out.append(md_table_row(k, v))
            out.append("")
    return "\n".join(out)

def required_attrs(block):
    return [k for k, v in (block.get("attributes") or {}).items() if v.get("required")]

def example_provider(provider_key, pblock):
    lines = [
        'terraform {',
        '  required_providers {',
        f'    autoglue = {{',
        f'      source  = "{provider_key}"',
        '      # version = ">= 0.0.0"',
        '    }',
        '  }',
        '}',
        '',
        'provider "autoglue" {',
    ]
    for k, v in (pblock.get("attributes") or {}).items():
        if v.get("required"):
            lines.append(f'  {k} = "REQUIRED_{k.upper()}"')
    for k, v in (pblock.get("attributes") or {}).items():
        if v.get("optional") and not v.get("computed"):
            lines.append(f'  # {k} = "..."')
    lines.append('}')
    return "```hcl\n" + "\n".join(lines) + "\n```"

def example_resource(rname, rblock):
    reqs = required_attrs(rblock)
    if reqs:
        body = "\n".join(f'  {k} = "..."' for k in reqs)
    else:
        body = "  # no required attributes"
    return f"```hcl\nresource \"{rname}\" \"example\" {{\n{body}\n}}\n```"

def example_data(dname, _dblock):
    return f"""```hcl
data "{dname}" "all" {{}}

# Example of reading exported fields (adjust to your needs):
# output "first_item_raw" {{
#   value = try(data.{dname}.all.items[0].raw, null)
# }}
```"""

def example_function_call(provider_local_name, fname):
    return f"""```hcl
# Example of calling a provider function (if available)
# local value = {provider_local_name}::{fname}("arg")
```"""

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--schema", required=True)
    ap.add_argument("--provider", required=True, help="provider key e.g. glueops/autoglue/autoglue")
    ap.add_argument("--out", required=True)
    args = ap.parse_args()

    with open(args.schema, "r", encoding="utf-8") as f:
        doc = json.load(f)

    prov = (doc.get("provider_schemas") or {}).get(args.provider)
    if not prov:
        sys.exit(f"Provider '{args.provider}' not found in schema.")

    out = []
    out.append(f"# {args.provider} – Reference (generated)\n")
    out.append("_Generated from providers schema JSON._\n")

    # Provider config
    pblock = (prov.get("provider") or {}).get("block", {}) or {}
    out.append("## Provider Configuration\n")
    out.append(render_block(pblock) or "_No provider configuration attributes found._")
    out.append("\n### Basic usage\n")
    out.append(example_provider(args.provider, pblock))

    # Functions
    funcs = prov.get("functions") or {}
    out.append("## Provider Functions\n")
    if not funcs:
        out.append("_No provider-defined functions._\n")
    else:
        for fname, fdef in funcs.items():
            out.append(f"### `{fname}`\n")
            out.append((fdef.get("summary") or fdef.get("description") or "").strip() + "\n")
            out.append(f"- **Return type:** `{type_to_str(fdef.get('return_type'))}`\n")
            params = fdef.get("parameters") or []
            if params:
                out.append("\n**Parameters**\n")
                out.append("| Name | Type | Description |\n|---|---|---|\n")
                for p in params:
                    pdesc = (p.get("description") or "").replace("\n", " ")
                    out.append(f"| `{p.get('name')}` | {type_to_str(p.get('type'))} | {pdesc} |\n")
            out.append("\n")
            out.append(example_function_call("autoglue", fname))

    # Resources
    resources = prov.get("resource_schemas") or {}
    out.append("## Resources\n")
    if not resources:
        out.append("_None._\n")
    else:
        for rname in sorted(resources.keys()):
            rs = resources[rname]
            rblock = rs.get("block", {}) or {}
            out.append(f"### `{rname}`\n")
            if rblock.get("description"):
                out.append(rblock["description"] + "\n")
            out.append(render_block(rblock) or "_No attributes._")
            out.append("\n**Example**\n")
            out.append(example_resource(rname, rblock))

    # Data sources
    datas = prov.get("data_source_schemas") or {}
    out.append("## Data Sources\n")
    if not datas:
        out.append("_None._\n")
    else:
        for dname in sorted(datas.keys()):
            ds = datas[dname]
            dblock = ds.get("block", {}) or {}
            out.append(f"### `{dname}`\n")
            if dblock.get("description"):
                out.append(dblock["description"] + "\n")
            out.append(render_block(dblock) or "_No attributes._")
            out.append("\n**Example**\n")
            out.append(example_data(dname, dblock))

    with open(args.out, "w", encoding="utf-8") as f:
        f.write("\n".join(out))

if __name__ == "__main__":
    main()
