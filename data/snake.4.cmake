# Copyright (c) 2022-2024 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

# Handle packaging via CPack
#
# TODO: This is a work in progress... Not ready for use.
# -------------------------------------------------------------------------------------------------------
if(SNAKE_ENABLE_PACKAGING)
  include(InstallRequiredSystemLibraries)

  set_opt(CPACK_PACKAGE_NAME "${CMAKE_PROJECT_NAME}")
  set_opt(CPACK_PACKAGE_DESCRIPTION_FILE "${CMAKE_SOURCE_DIR}/README.md")
  set_opt(CPACK_RESOURCE_FILE_LICENSE "${CMAKE_SOURCE_DIR}/LICENSE")
  set_opt(CPACK_PACKAGE_INSTALL_DIRECTORY "${CMAKE_PROJECT_NAME}")
  set_opt(CPACK_PACKAGE_CONTACT "${SNAKE_CONTACT}")
  set_opt(CPACK_PACKAGE_VENDOR "${SNAKE_ORGANIZATION}")
  set_opt(CPACK_PACKAGE_DESCRIPTION_SUMMARY "${CMAKE_PROJECT_DESCRIPTION}")
  set_opt(CPACK_PACKAGE_EXECUTABLES "true")
  set_opt(CPACK_CREATE_DESKTOP_LINKS "true")
  set_opt(CPACK_STRIP_FILES "true")

  if(PLATFORM_LINUX)
    set(RPM_ENABLED OFF)
    set(DEB_ENABLED ON)
    set(TXZ_ENABLED ON)

    find_program(RPMBUILD_PATH rpmbuild)
    if(RPMBUILD_PATH)
      set(RPM_ENABLED ON)
    endif()

    if(RPM_ENABLED)
      list(APPEND CPACK_GENERATOR "RPM")
      set(CPACK_RPM_PACKAGE_LICENSE "${SNAKE_PROJECT_LICENSE}")
    endif()

    if(DEB_ENABLED)
      list(APPEND CPACK_GENERATOR "DEB")
      set(CPACK_DEBIAN_PACKAGE_CONTROL_STRICT_PERMISSION TRUE)
      set(CPACK_DEBIAN_PACKAGE_ARCHITECTURE "amd64")
      set(CPACK_DEBIAN_PACKAGE_HOMEPAGE "${CMAKE_PROJECT_HOMEPAGE_URL}")
      set(CPACK_DEBIAN_COMPRESSION_TYPE "zstd")
    endif()

    if(TXZ_ENABLED)
      list(APPEND CPACK_GENERATOR "TXZ")
    endif()

    install(FILES "${CMAKE_SOURCE_DIR}/LICENSE" DESTINATION "share/licenses/${SNAKE_PROJECT_NAME}")
  elseif(PLATFORM_OSX)
    # Linux CPack packages.
    # ------------------------------------------------------------------------------
    set(BUNDLE_ENABLED ON)

    if(BUNDLE_ENABLED)
      list(APPEND CPACK_GENERATOR "DragNDrop")
      set(CPACK_PACKAGING_INSTALL_PREFIX "/")
      set(CPACK_OSX_PACKAGE_VERSION "10.6")
      set(MACOSX_BUNDLE_INFO_PLIST "${OSX_PACKAGE_SOURCE_DIR}/info.plist")
    endif()
  elseif(PLATFORM_ANDROID)
    # Android packages.
  elseif(PLATFORM_IOS)
    # iOS packages.
  elseif(PLATFORM_WINDOWS)
    set(ZIP_ENABLE ON)
    set(NSIS_ENABLE OFF)

    find_program(NSIS_PATH nsis PATH_SUFFIXES nsis)
    if(NSIS_PATH)
      set(NSIS_ENABLE ON)
    endif()

    if(ZIP_ENABLE)
      list(APPEND CPACK_GENERATOR "ZIP")
    endif()

    if(NSIS_ENABLE)
      list(APPEND CPACK_GENERATOR "NSIS")
      set(CPACK_NSIS_DISPLAY_NAME "${CPACK_PACKAGE_NAME}")
      set(CPACK_NSIS_MUI_ICON "${WINDOWS_PACKAGE_SOURCE_DIR}/icon.ico")
      set(CPACK_NSIS_HELP_LINK "${CMAKE_PROJECT_HOMEPAGE_URL}")
      set(CPACK_NSIS_URL_INFO_ABOUT "${CMAKE_PROJECT_HOMEPAGE_URL}")
      set(CPACK_NSIS_CONTACT "${CPACK_PACKAGE_CONTACT}")
      set(CPACK_NSIS_MODIFY_PATH ON)
    endif()
  endif()

  include(CPack)
endif()

# Symbolically link the compilation database to the directory above
# -------------------------------------------------------------------------------------------------------
if(CMAKE_EXPORT_COMPILE_COMMANDS)
  file(CREATE_LINK "${CMAKE_BINARY_DIR}/compile_commands.json" "${SNAKE_DIR}/compile_commands.json" SYMBOLIC)
endif()

# Print variables
# -------------------------------------------------------------------------------------------------------
print_status("Snake Variables")

print_variable("SNAKE_ENABLE_DEVELOPMENT_MODE")
print_variable("SNAKE_ENABLE_TESTING")
print_variable("SNAKE_ENABLE_TRANSLATIONS")
print_variable("SNAKE_ENABLE_VERBOSE")
print_variable("SNAKE_ENABLE_PACKAGING")
print_variable("SNAKE_ENABLE_DOCUMENTATION")
print_variable("SNAKE_PROJECT_LIBRARY_DIRECTORY")
print_variable("SNAKE_PROJECT_i18n_DIRECTORY")
print_variable("SNAKE_PROJECT_MODULE_DIRECTORY")
print_variable("SNAKE_PROJECT_BINARY_DIRECTORY")
print_variable("SNAKE_ENABLE_QML_MODULES")
print_variable("SNAKE_ORGANIZATION")
print_variable("SNAKE_ORGANIZATION_NAME")
print_variable("SNAKE_PROJECT_NAME")
print_variable("SNAKE_DIR")

print_status("CMake Variables")

print_variable("CMAKE_EXPORT_COMPILE_COMMANDS")
print_variable("CMAKE_BUILD_TYPE")
print_variable("CMAKE_POSITION_INDEPENDENT_CODE")
print_variable("CMAKE_AUTOGEN_VERBOSE")
print_variable("CMAKE_AUTOMOC_PATH_PREFIX")
print_variable("CMAKE_ENABLE_EXPORTS")
print_variable("CMAKE_UNITY_BUILD")
print_variable("CMAKE_CONFIGURATION_TYPES")
print_variable("CMAKE_ERROR_DEPRECATED")
print_variable("CMAKE_OPTIMIZE_DEPENDENCIES")
print_variable("CMAKE_CXX_FLAGS")
print_variable("CMAKE_CXX_COMPILER_LAUNCHER")
print_variable("CMAKE_CXX_EXTENSIONS")
print_variable("CMAKE_CXX_STANDARD")
print_variable("CMAKE_CXX_STANDARD_REQUIRED")
print_property("COMPILE_OPTIONS")
print_property("LINK_OPTIONS")
print_variable("CMAKE_ENABLE_CLANG_TIDY")
print_variable("CMAKE_CXX_COMPILER")
print_variable("CMAKE_INTERPROCEDURAL_OPTIMIZATION")
print_variable("CMAKE_COLOR_MAKEFILE")
print_variable("CMAKE_STRIP_OUTPUT_BINARIES")
print_variable("CMAKE_SOURCE_DIR")
print_variable("CMAKE_INSTALL_PREFIX")
print_variable("CMAKE_BINARY_DIR")
print_property("PACKAGES_FOUND")
print_property("PACKAGES_NOT_FOUND")
print_variable("CMAKE_PROJECT_NAME")
print_variable("CMAKE_PROJECT_VERSION")
print_variable("CMAKE_SYSTEM")
print_variable("CMAKE_SYSTEM_PROCESSOR")
print_variable("CMAKE_MAKE_PROGRAM")
print_variable("CMAKE_VERSION")

# Formatting target.
# -------------------------------------------------------------------------------------------------------
find_program(CLANG_FORMAT_PROGRAM clang-format)
if(CLANG_FORMAT_PROGRAM)
  print_message(STATUS "Enabling clang-format...")
  add_custom_target(format COMMAND ${CLANG_FORMAT_PROGRAM} -i ${ALL_CXX_SOURCES} VERBATIM)
endif()
