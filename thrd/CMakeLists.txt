#######################################################################
#                   THIRD PARTY LIB CMAKE
#######################################################################
#Use fmtlib for formatting
add_subdirectory(fmt EXCLUDE_FROM_ALL)
target_set_warnings(fmt DISABLE ALL)
set_target_properties(fmt PROPERTIES FOLDER "external")


#Use spdlog for logging, and tell it to use our version of fmtlib
add_subdirectory(spdlog EXCLUDE_FROM_ALL)
target_compile_definitions(spdlog INTERFACE SPDLOG_FMT_EXTERNAL)
file(MAKE_DIRECTORY "${CMAKE_CURRENT_BINARY_DIR}/logs")


# Doctest for unit tests
add_library(doctest INTERFACE)
target_include_directories(
    doctest
    INTERFACE
    doctest/doctest # note : will expose the parts/ folder...
)
target_compile_features(doctest INTERFACE cxx_std_11)
include_directories(${DOCTEST_INCLUDE_DIR})

#argument passing
add_library(args INTERFACE)
target_include_directories(
    args
    INTERFACE
    args # note : will expose the parts/ folder...
)

#colored output
add_library(rang INTERFACE)
target_include_directories(
    rang
    INTERFACE
    rang/include # note : will expose the parts/ folder...
)

#YAML parsing
add_library(yaml-cpp INTERFACE)
target_include_directories(
    yaml-cpp
    INTERFACE
    yaml-cpp/include/yaml-cpp # note : will expose the parts/ folder...
)
include_directories(yaml-cpp${CMAKE_STATIC_LIBRARY_SUFFIX})
