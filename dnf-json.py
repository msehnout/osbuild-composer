#!/usr/bin/python3

import datetime
import dnf
import hashlib
import hawkey
import json
import shutil
import sys
import tempfile

DNF_ERROR_EXIT_CODE = 10


def timestamp_to_rfc3339(timestamp):
    d = datetime.datetime.utcfromtimestamp(package.buildtime)
    return d.strftime('%Y-%m-%dT%H:%M:%SZ')


def dnfrepo(desc, parent_conf=None):
    """Makes a dnf.repo.Repo out of a JSON repository description"""

    repo = dnf.repo.Repo(desc["id"], parent_conf)

    if "baseurl" in desc:
        repo.baseurl = desc["baseurl"]
    elif "metalink" in desc:
        repo.metalink = desc["metalink"]
    elif "mirrorlist" in desc:
        repo.mirrorlist = desc["mirrorlist"]
    else:
        assert False

    if desc.get("ignoressl", False):
        repo.sslverify = False

    return repo


def create_base(repos, module_platform_id, persistdir, cachedir):
    base = dnf.Base()
    base.conf.module_platform_id = module_platform_id
    base.conf.config_file_path = "/dev/null"
    base.conf.persistdir = persistdir
    base.conf.cachedir = cachedir

    for repo in repos:
        base.repos.add(dnfrepo(repo, base.conf))

    base.fill_sack(load_system_repo=False)
    return base


def exit_with_dnf_error(kind: str, reason: str):
    json.dump({"kind": kind, "reason": reason}, sys.stdout)
    sys.exit(DNF_ERROR_EXIT_CODE)


def repo_checksums(base):
    checksums = {}
    for repo in base.repos.iter_enabled():
        # Uses the same algorithm as libdnf to find cache dir:
        #   https://github.com/rpm-software-management/libdnf/blob/master/libdnf/repo/Repo.cpp#L1288
        if repo.metalink:
            url = repo.metalink
        elif repo.mirrorlist:
            url = repo.mirrorlist
        elif repo.baseurl:
            url = repo.baseurl[0]
        else:
            assert False

        digest = hashlib.sha256(url.encode()).hexdigest()[:16]

        with open(f"{base.conf.cachedir}/{repo.id}-{digest}/repodata/repomd.xml", "rb") as f:
            repomd = f.read()

        checksums[repo.id] = "sha256:" + hashlib.sha256(repomd).hexdigest()

    return checksums


call = json.load(sys.stdin)
command = call["command"]
arguments = call["arguments"]
repos = arguments.get("repos", {})
cachedir = arguments["cachedir"]
module_platform_id = arguments["module_platform_id"]

with tempfile.TemporaryDirectory() as persistdir:
    try:
        base = create_base(repos, module_platform_id, persistdir, cachedir)
    except dnf.exceptions.RepoError as e:
        exit_with_dnf_error("RepoError", f"Error occurred when setting up repo: {e}")

    if command == "dump":
        packages = []
        for package in base.sack.query().available():
            packages.append({
                "name": package.name,
                "summary": package.summary,
                "description": package.description,
                "url": package.url,
                "epoch": package.epoch,
                "version": package.version,
                "release": package.release,
                "arch": package.arch,
                "buildtime": timestamp_to_rfc3339(package.buildtime),
                "license": package.license
            })
        json.dump({
            "checksums": repo_checksums(base),
            "packages": packages
        }, sys.stdout)

    elif command == "depsolve":
        errors = []

        try:
            base.install_specs(arguments["package-specs"], exclude=arguments.get("exclude-specs", []))
        except dnf.exceptions.MarkingErrors as e:
            exit_with_dnf_error("MarkingErrors", f"Error occurred when marking packages for installation: {e}")

        try:
            base.resolve()
        except dnf.exceptions.DepsolveError as e:
            exit_with_dnf_error("DepsolveError", f"There was a problem depsolving {arguments['package-specs']}: {e}")

        dependencies = []
        for tsi in base.transaction:
            # avoid using the install_set() helper, as it does not guarantee a stable order
            if tsi.action not in dnf.transaction.FORWARD_ACTIONS:
                continue
            package = tsi.pkg

            dependencies.append({
                "name": package.name,
                "epoch": package.epoch,
                "version": package.version,
                "release": package.release,
                "arch": package.arch,
                "repo_id": package.reponame,
                "path": package.relativepath,
                "remote_location": package.remote_location(),
                "checksum": f"{hawkey.chksum_name(package.chksum[0])}:{package.chksum[1].hex()}",
            })
        json.dump({
            "checksums": repo_checksums(base),
            "dependencies": dependencies
        }, sys.stdout)