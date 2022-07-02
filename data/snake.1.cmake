# Copyright (c) 2022 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

# Disable in-source builds.
if(CMAKE_SOURCE_DIR STREQUAL CMAKE_BINARY_DIR)
  message(FATAL_ERROR "In-source builds are not allowed.")
endif()

# Only sets the option if it has not been defined (by a user or profile).
macro(set_opt K V)
  if(NOT DEFINED ${K})
    set(${K} ${V})
  endif()
endmacro()

# Check to see if snake.lock has been deleted.
if(NOT EXISTS "${SNAKE_DIR}/snake.lock")
  set(SNAKE_FORCE_UPDATE on)
  file(TOUCH "${SNAKE_DIR}/snake.lock")
else()
  set(SNAKE_FORCE_UPDATE off)
endif()

# Memoization for performance.
# -------------------------------------------------------------------------------------------------------
macro(find_package)
  if(NOT TARGET ${ARGV0})
    set(ARG_LIST ${ARGN})
    list(JOIN ARG_LIST "" CACHE_KEY)
    if(NOT DEFINED _find_package${CACHE_KEY})
      set(_find_package${CACHE_KEY} on)
      _find_package(${ARGN})
    endif()
  endif()
endmacro()

macro(find_library)
  set(ARG_LIST ${ARGN})
  list(JOIN ARG_LIST "" CACHE_KEY)
  if(NOT DEFINED _find_library${CACHE_KEY})
    set(_find_library${CACHE_KEY} on)
    _find_library(${ARGN})
  endif()
endmacro()

macro(include)
  set(ARG_LIST ${ARGN})
  list(JOIN ARG_LIST "" CACHE_KEY)
  if(NOT DEFINED _include${CACHE_KEY})
    set(_include${CACHE_KEY} on)
    _include(${ARGN})
  endif()
endmacro()

macro(add_subdirectory)
  set(ARG_LIST ${ARGN})
  list(JOIN ARG_LIST "" CACHE_KEY)
  if(NOT DEFINED _add_subdirectory${CACHE_KEY})
    set(_add_subdirectory${CACHE_KEY} on)
    _add_subdirectory(${ARGN})
  endif()
endmacro()

macro(find_program)
  set(ARG_LIST ${ARGN})
  list(JOIN ARG_LIST "" CACHE_KEY)
  if(NOT DEFINED _find_program${CACHE_KEY})
    set(_find_program${CACHE_KEY} on)
    _find_program(${ARGN})
  endif()
endmacro()

# Snake options
option(SNAKE_ENABLE_TESTING "Enable CMake testing (add this to 'requirement:')" off)
option(SNAKE_ENABLE_TRANSLATIONS "Enable Qt translations" off)
option(SNAKE_ENABLE_PACKAGING "Enable CMake packaging" off)
option(SNAKE_ENABLE_DEVELOPMENT_MODE "Enable development mode" off)
option(SNAKE_ENABLE_VERBOSE "Enable verbose CMake configuration" off)
option(SNAKE_ENABLE_QML_MODULES "Enable Qt QML modules" off)
option(SNAKE_ENABLE_DOCUMENTATION "Enable documentation generation" off)
option(SNAKE_ENABLE_COLOR "Enable colored configuration and compilation output" on)
option(SNAKE_ENABLE_EXAMPLES "Enable examples (add this to 'requirement:')" off)

option(SNAKE_ENABLE_EXPORT_PREFIX "Enable a project prefix for the export header" on)

option(SNAKE_DIR "The snake directory" "")
option(NO_SNAKE "Set this to true when you are not using the Snake executable" on)
option(SNAKE_ORGANIZATION "The snake organization" "")

# Silence Conan logs.
set(CONAN_CMAKE_SILENT_OUTPUT on)

# Calculate timestamp.
string(TIMESTAMP PROJECT_TIMESTAMP "%s")
set_opt(PROJECT_TIMESTAMP ${PROJECT_TIMESTAMP})

# The template directory (.h.in).
set(SNAKE_TEMPLATE_DIR "${SNAKE_DIR}/.templates")

# The generated include directory.
set(SNAKE_GENERATED_INCLUDE_DIR "${CMAKE_BINARY_DIR}/include")

# This variable is always set to 'on'.
set(SNAKE_ALWAYS_BUILD on)

# Disable C/CXX compiler extensions. These extensions can cause problems in cross-platform
# builds because they aren't compatible between compilers.
set_opt(CMAKE_C_EXTENSIONS OFF)
set_opt(CMAKE_CXX_EXTENSIONS OFF)

# Set the default standard to C++20 (can be changed/reduced).
set_opt(CMAKE_CXX_STANDARD 20)

# Symbols are hidden by default.
set_opt(CMAKE_CXX_VISIBILITY_PRESET hidden)
set_opt(CMAKE_VISIBILITY_INLINES_HIDDEN on)

# Enable colors.
if(SNAKE_ENABLE_COLOR)
  set(CMAKE_COLOR_MAKEFILE on)
else()
  set(CMAKE_COLOR_MAKEFILE off)
endif()

# The default CMake module path. Should not be changed.
set(CMAKE_MODULE_PATH "${SNAKE_DIR}/.cmake" CACHE INTERNAL "")

list(APPEND CMAKE_MODULE_PATH ${CMAKE_BINARY_DIR})
list(APPEND CMAKE_PREFIX_PATH ${CMAKE_BINARY_DIR})

# Find the snake binary.
find_program(SNAKE_PROGRAM snake)

# Use ccache by default for C++ programs if found.
find_program(CCACHE_PROGRAM ccache)
if(CCACHE_PROGRAM)
  set_opt(CMAKE_CXX_COMPILER_LAUNCHER "${CCACHE_PROGRAM}")
endif()

# Build shared libraries by default.
set_opt(BUILD_SHARED_LIBS on)

if(BUILD_SHARED_LIBS)
  set(SNAKE_LIB_TYPE SHARED)
else()
  set(SNAKE_LIB_TYPE STATIC)
endif()

# Export commands by default.
set(CMAKE_EXPORT_COMPILE_COMMANDS on CACHE INTERNAL "")

# Require the C++ standard specified by the user.
set_opt(CMAKE_CXX_STANDARD_REQUIRED ON)

# Lowercase project name
string(TOLOWER "${CMAKE_PROJECT_NAME}" SNAKE_PROJECT_NAME)

# Lowercase organization name.
string(TOLOWER "${SNAKE_ORGANIZATION}" SNAKE_ORGANIZATION_NAME)

# Default install directory
if(CMAKE_INSTALL_PREFIX_INITIALIZED_TO_DEFAULT)
  set_opt(CMAKE_INSTALL_PREFIX "${CMAKE_BINARY_DIR}/install")
endif()

# Generic locations
set_opt(SNAKE_LIBRARY_DIRECTORY "lib/${SNAKE_PROJECT_NAME}/dynamic")
set_opt(SNAKE_ARCHIVE_DIRECTORY "lib/${SNAKE_PROJECT_NAME}/archive")
set_opt(SNAKE_MODULES_DIRECTORY "lib/${SNAKE_PROJECT_NAME}/modules")
set_opt(SNAKE_RUNTIME_DIRECTORY "bin")
set_opt(SNAKE_TRANSLATIONS_DIR "share/${SNAKE_PROJECT_NAME}/i18n")

# Output directories
set_opt(CMAKE_RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/${SNAKE_RUNTIME_DIRECTORY}")
set_opt(CMAKE_LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/${SNAKE_LIBRARY_DIRECTORY}")
set_opt(CMAKE_ARCHIVE_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/${SNAKE_ARCHIVE_DIRECTORY}")
set_opt(CMAKE_MODULES_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/${SNAKE_MODULES_DIRECTORY}")
set_opt(QT_QML_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/qml")

if(SNAKE_ENABLE_DEVELOPMENT_MODE)
  # Development directories
  set_opt(SNAKE_PROJECT_LIBRARY_DIRECTORY "${CMAKE_LIBRARY_OUTPUT_DIRECTORY}")
  set_opt(SNAKE_PROJECT_ARCHIVE_DIRECTORY "${CMAKE_ARCHIVE_OUTPUT_DIRECTORY}")
  set_opt(SNAKE_PROJECT_i18n_DIRECTORY "${CMAKE_BINARY_DIR}/${TRANSLATIONS_DIR}")
  set_opt(SNAKE_PROJECT_MODULE_DIRECTORY "${CMAKE_MODULES_OUTPUT_DIRECTORY}")
  set_opt(SNAKE_PROJECT_BINARY_DIRECTORY "${CMAKE_RUNTIME_OUTPUT_DIRECTORY}")
else()
  # Installation directories
  set_opt(SNAKE_PROJECT_LIBRARY_DIRECTORY "${CMAKE_INSTALL_PREFIX}/${SNAKE_LIBRARY_DIRECTORY}")
  set_opt(SNAKE_PROJECT_ARCHIVE_DIRECTORY "${CMAKE_INSTALL_PREFIX}/${SNAKE_ARCHIVE_DIRECTORY}")
  set_opt(SNAKE_PROJECT_i18n_DIRECTORY "${CMAKE_INSTALL_PREFIX}/${TRANSLATIONS_DIR}")
  set_opt(SNAKE_PROJECT_MODULE_DIRECTORY "${CMAKE_INSTALL_PREFIX}/${SNAKE_MODULES_DIRECTORY}")
  set_opt(SNAKE_PROJECT_BINARY_DIRECTORY "${CMAKE_INSTALL_PREFIX}/${SNAKE_RUNTIME_DIRECTORY}")
endif()

# Add some useful platform variables and definitions
# -------------------------------------------------------------------------------------------------------
if(${CMAKE_SYSTEM_NAME} STREQUAL "Linux")
  set(LINUX on)
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "iOS")
  set(IOS on)
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Android")
  set(ANDROID on)
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Darwin")
  set(OSX on)
elseif(${CMAKE_SYSTEM_NAME} STREQUAL "Windows")
  set(WINDOWS on)
endif()

if(LINUX OR IOS OR OSX OR ANDROID)
  set(UNIX on)
endif()

if(LINUX OR WINDOWS OR OSX)
  set(DESKTOP on)
endif()

if(ANDROID OR IOS)
  set(MOBILE on)
  add_compile_definitions(Q_OS_MOBILE)
endif()

# Enable testing
# -------------------------------------------------------------------------------------------------------
if(SNAKE_ENABLE_TESTING)
  enable_testing()
endif()

# Add compilation options
# -------------------------------------------------------------------------------------------------------
add_compile_options(${SNAKE_GLOBAL_COMPILE_OPTIONS})
add_link_options(${SNAKE_GLOBAL_LINKER_OPTIONS})

# Enable IPO
# -------------------------------------------------------------------------------------------------------
if(CMAKE_INTERPROCEDURAL_OPTIMIZATION)
  include(CheckIPOSupported)
  check_ipo_supported(RESULT IPO_SUPPORTED)
  if(NOT IPO_SUPPORTED)
    set(CMAKE_INTERPROCEDURAL_OPTIMIZATION off)
  endif()
endif()

# Enable color for generators
# -------------------------------------------------------------------------------------------------------
if(SNAKE_ENABLE_COLOR)
  if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU")
    add_compile_options("-fdiagnostics-color=always")
  elseif(CMAKE_CXX_COMPILER_ID STREQUAL "Clang")
    add_compile_options("-fcolor-diagnostics")
    add_compile_options("-fno-diagnostics-show-option")
  endif()
endif()

# Define colors
# ------------------------------------------------------------------------------
if(SNAKE_ENABLE_COLOR)
  string(ASCII 27 Esc)
  set(COLOR_RESET "${Esc}[m")
  set(COLOR_BOLD "${Esc}[1m")
  set(COLOR_RED "${Esc}[31m")
  set(COLOR_GREEN "${Esc}[32m")
  set(COLOR_BLUE "${Esc}[34m")
  set(COLOR_CYAN "${Esc}[36m")
  set(COLOR_MAGENTA "${Esc}[35m")
  set(COLOR_YELLOW "${Esc}[33m")
  set(COLOR_WHITE "${Esc}[37m")
  set(COLOR_BOLD_RED "${Esc}[1;31m")
  set(COLOR_BOLD_GREEN "${Esc}[1;32m")
  set(COLOR_BOLD_BLUE "${Esc}[1;34m")
  set(COLOR_BOLD_CYAN "${Esc}[1;36m")
  set(COLOR_BOLD_MAGENTA "${Esc}[1;35m")
  set(COLOR_BOLD_YELLOW "${Esc}[1;33m")
  set(COLOR_BOLD_WHITE "${Esc}[1;37m")
  set(COLOR_GRAY "${Esc}[1;30m")
endif()
