##############################################################################################################################################
#                 CMAKE
##############################################################################################################################################
###### HUNTER PACKAGE MANAGER ######
include(ExternalProject)
find_package(Git REQUIRED)
include("cmake/HunterGate.cmake")
HunterGate(
    URL "https://github.com/ruslo/hunter/archive/v0.23.214.tar.gz"
    SHA1 "e14bc153a7f16d6a5eeec845fb0283c8fad8c358"
)

###### GENERAL SETTINGS ######
cmake_minimum_required(VERSION 3.10)
set(CMAKE_EXPORT_COMPILE_COMMANDS 1)
set(PROJECT_NAME "CLIntacts")
project(${PROJECT_NAME})
message ("\nPROJECT: ${CMAKE_PROJECT_NAME}")


######  VARIABLES  ######
set(${CMAKE_PROJECT_NAME}_SUMMARY               "A simple CONTACTS for CLI app")
set(${CMAKE_PROJECT_NAME}_VENDOR_NAME           "bresilla")
set(${CMAKE_PROJECT_NAME}_VENDOR_CONTACT        "trim.bresilla@gmail.com")
set(${CMAKE_PROJECT_NAME}_VERSION_MAJOR         0)
set(${CMAKE_PROJECT_NAME}_VERSION_MINOR         1)
set(${CMAKE_PROJECT_NAME}_VERSION_PATCH         0)
set(PROJECT_VERSION                             "${${CMAKE_PROJECT_NAME}_VERSION_MAJOR}.${${CMAKE_PROJECT_NAME}_VERSION_MINOR}.${${CMAKE_PROJECT_NAME}_VERSION_PATCH}")
set(CMAKE_MODULE_PATH                          "${CMAKE_CURRENT_SOURCE_DIR}/cmake")
set(PROJECT_PATH                               "${CMAKE_CURRENT_SOURCE_DIR}")
set(PROJECT_BUILD                              "${PROJECT_PATH}/build")
set(MAIN_SRC                                   "${PROJECT_PATH}/src")
set(TEST_SRC                                   "${PROJECT_PATH}/test")
set(THRD_SRC                                   "${PROJECT_PATH}/thrd")
set(LIB_BUILD                                  "${PROJECT_BUILD}/lib")
set(BIN_BUILD                                  "${PROJECT_BUILD}/bin")
set(INCLUDE_BUILD                              "${PROJECT_BUILD}/include")
set(DOCS_BUILD                                 "${PROJECT_BUILD}/docs")
set(EXECUTABLE_OUTPUT_PATH                     ${BIN_BUILD})
set(LIBRARY_OUTPUT_PATH                        ${LIB_BUILD})


####### BUILD TYPE #######
string(TOLOWER "${CMAKE_BUILD_TYPE}" cmake_build_type_tolower)
if(NOT cmake_build_type_tolower STREQUAL "debug" AND NOT cmake_build_type_tolower STREQUAL "release" AND NOT cmake_build_type_tolower STREQUAL "fastdebug" AND NOT cmake_build_type_tolower STREQUAL "fastrelease")
    message(FATAL_ERROR "Unknown build type \"${CMAKE_BUILD_TYPE}\". Allowed values are Debug, Release, FastDebug, and FastRelease (case-insensitive).")
endif()
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE "Release")
endif()


####### COMPILER ######
include(Warnings)
if((cmake_build_type_tolower STREQUAL "release") OR (cmake_build_type_tolower STREQUAL "fastrelease"))
    set(CMAKE_CXX_FLAGS         "${CMAKE_CXX_FLAGS} -Wall -Wno-unknown-pragmas -Wno-sign-compare -Woverloaded-virtual -Wwrite-strings -Wno-unused")
    set(CMAKE_CXX_FLAGS_RELEASE "-O3")
    include(LTO)
    include(Packaging)
    option(BUILD_MAIN                              "Build main function"            ON)
    if(cmake_build_type_tolower STREQUAL "release")
        option(ENABLE_DOCUMENT                         "Enable documentation"           ON)
    endif()
elseif((cmake_build_type_tolower STREQUAL "debug") OR (cmake_build_type_tolower STREQUAL "fastdebug"))
    set(CMAKE_CXX_FLAGS         "${CMAKE_CXX_FLAGS} -Wall -Wno-unknown-pragmas -Wno-sign-compare -Woverloaded-virtual -Wwrite-strings -Wno-unused")
    set(CMAKE_CXX_FLAGS         "${CMAKE_CXX_FLAGS} --coverage -fprofile-arcs -ftest-coverage")
    set(CMAKE_CXX_FLAGS_DEBUG   "-O0 -g3")
    option(BUILD_MAIN                              "Build main function"            ON)
    if(cmake_build_type_tolower STREQUAL "debug")
        option(BUILD_TESTS                             "Build test units"               ON)
        option(ENABLE_COVERAGE                         "Enable testing coverage"        ON)
        option(ENABLE_SANITIZERS                       "Enable langtidy and cppcheck"   ON)
        option(ENABLE_DOCUMENT                         "Enable documentation"           ON)
        option(ENABLE_MEMCHECK                         "Enable memory leak check"       ON)
    endif()
endif()
message ("BUILD TYPE: ${CMAKE_BUILD_TYPE}\n")
add_subdirectory(${PROJECT_SOURCE_DIR}/thrd EXCLUDE_FROM_ALL)
add_subdirectory(${PROJECT_SOURCE_DIR}/src)


#CPP STANDARD OPTIONS
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
