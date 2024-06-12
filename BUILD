load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "thing",
    srcs = ["main.go"],
    importpath = "example.com/my/thing",
    visibility = ["//visibility:public"],
    deps = ["@com_github_mattn_go_ieproxy//:go_default_library"],
)
