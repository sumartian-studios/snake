# Copyright (c) 2022-2024 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

# C++20 standard features...
# -------------------------------------------------------------------------------------------------------
if(CMAKE_CXX_STANDARD_REQUIRED AND (CMAKE_CXX_STANDARD GREATER_EQUAL 20))
  if(CMAKE_CXX_COMPILER_ID STREQUAL "Clang" OR CMAKE_CXX_COMPILER_ID STREQUAL "AppleClang")
    add_compile_options(-fcoroutines-ts)
  endif()
endif()

# Download/install dependencies...
# -------------------------------------------------------------------------------------------------------
if(SNAKE_FORCE_UPDATE)
  message(STATUS "Snake installing dependencies...")

  # conan_cmake_configure(REQUIRES ${ENABLED_CONAN_PACKAGES} GENERATORS "cmake_find_package")
  # conan_cmake_autodetect(settings)
  # conan_cmake_install(PATH_OR_REFERENCE
  #                     "${CMAKE_BINARY_DIR}"
  #                     BUILD
  #                     "missing"
  #                     REMOTE
  #                     "conancenter"
  #                     SETTINGS
  #                     "${settings}"
  #                     OUTPUT_QUIET
  # )
endif()
