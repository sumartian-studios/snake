# Copyright (c) 2022-2024 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

# Print a top-level status.
# -------------------------------------------------------------------------------------------------------
macro(print_status)
  if(SNAKE_ENABLE_VERBOSE)
    message("${COLOR_BOLD_WHITE}-- ${ARGN}${COLOR_RESET}")
  endif()
endmacro()

macro(print_dim_status)
  if(SNAKE_ENABLE_VERBOSE)
    message("${COLOR_GRAY}-- ${ARGN}${COLOR_RESET}")
  endif()
endmacro()

# Print a verbose message.
# -------------------------------------------------------------------------------------------------------
macro(print_message)
  if(SNAKE_ENABLE_VERBOSE)
    message(${ARGN})
  endif()
endmacro()

# Print a property.
# -------------------------------------------------------------------------------------------------------
macro(print_property)
  if(SNAKE_ENABLE_VERBOSE)
    get_property(PROPERTY GLOBAL PROPERTY ${ARGN})

    if(NOT PROPERTY)
      get_property(PROPERTY DIRECTORY . PROPERTY ${ARGN})
    endif()

    print_value("  " "${ARGN}" "${PROPERTY}")
  endif()
endmacro()

# Print a variable.
# -------------------------------------------------------------------------------------------------------
macro(print_variable)
  if(SNAKE_ENABLE_VERBOSE)
    print_value("  " "${ARGV0}" "${${ARGV0}}")
  endif()
endmacro()

# Print a value (using color)
# -------------------------------------------------------------------------------------------------------
macro(print_value)
  if(SNAKE_ENABLE_VERBOSE)
    # Colorize the arguments.
    if("${ARGV2}" STREQUAL "on" OR "${ARGV2}" STREQUAL "true")
      set(TITLE "${COLOR_GREEN}${ARGV1}:${COLOR_RESET}")
      set(VALUE " ${COLOR_BOLD_GREEN}${ARGV2}${COLOR_RESET}")
    elseif("${ARGV2}" STREQUAL "off" OR "${ARGV2}" STREQUAL "false" OR "${ARGV2}" STREQUAL "0")
      set(TITLE "${COLOR_RED}${ARGV1}:${COLOR_RESET}")
      set(VALUE " ${COLOR_BOLD_RED}${ARGV2}${COLOR_RESET}")
    elseif("${ARGV2}" STREQUAL "")
      set(TITLE "${COLOR_RED}${ARGV1}:${COLOR_RESET}")
      set(VALUE " ${COLOR_BOLD_RED}off${COLOR_RESET}")
    else()
      set(TITLE "${COLOR_GREEN}${ARGV1}:${COLOR_RESET}")

      # Replace ';' with ' ' in our value string. We only do this here because
      # the value can be a list.
      string(REPLACE ";" " " TEXT "${ARGV2}")

      set(VALUE " ${TEXT}")
    endif()

    # Print the formatted message.
    message("${ARGV0} ${TITLE}${VALUE}")
  endif()
endmacro()

# Add Qt resources
# -------------------------------------------------------------------------------------------------------
macro(snake_add_resources TARGET REGEXES MODULE PREFIX)
  unset(HAS_PREFIX)

  if(NOT ${PREFIX} STREQUAL "")
    set(HAS_PREFIX on)
  endif()

  foreach(REGEX ${REGEXES})
    file(GLOB_RECURSE RESOURCE_FILES ${REGEX})

    cmake_path(GET REGEX PARENT_PATH RELATIVE_SOURCE_DIR)

    foreach(FILE ${RESOURCE_FILES})
      # Get the file extension
      cmake_path(GET FILE EXTENSION FILE_EXT)

      # Get the file name
      cmake_path(GET FILE FILENAME RESOURCE_FILE_NAME)

      if(FILE_EXT STREQUAL ".qml")
        if(SNAKE_ENABLE_QML_MODULES)
          list(APPEND TARGET_QML_FILES "${FILE}")
        else()
          list(APPEND TARGET_QT_RESOURCE_FILES "${FILE}")
        endif()
      elseif(FILE_EXT STREQUAL ".ts" AND SNAKE_ENABLE_TRANSLATIONS)
        set_source_files_properties(${FILE} PROPERTIES OUTPUT_LOCATION "${CMAKE_BINARY_DIR}/${SNAKE_TRANSLATIONS_DIR}")
        list(APPEND TARGET_QT_TRANSLATION_FILES "${FILE}")
      elseif(FILE_EXT STREQUAL ".json" OR FILE_EXT STREQUAL ".yml")
        # Remove existing extension.
        cmake_path(REPLACE_EXTENSION RESOURCE_FILE_NAME ".cbor")

        # Set the path to the generated file.
        set_opt(OUTPUT_FILE "${CMAKE_BINARY_DIR}/.rcc/${TARGET}/${RESOURCE_FILE_NAME}")

        if(FILE_EXT STREQUAL ".json")
          set(MUTATOR "json-to-cbor")
        elseif(FILE_EXT STREQUAL ".yml")
          set(MUTATOR "yaml-to-cbor")
        endif()

        # Use the 'rccompiler' tool to generate output files.
        add_custom_command(OUTPUT "${OUTPUT_FILE}"
                           COMMAND ${SNAKE_PROGRAM} mutate -m ${MUTATOR} ${FILE} ${OUTPUT_FILE}
                           DEPENDS ${FILE}
                           VERBATIM
        )

        set(FILE ${OUTPUT_FILE})
        list(APPEND TARGET_QT_RESOURCE_FILES "${OUTPUT_FILE}")
        set(RESOURCE_GENERATED true)
      else()
        list(APPEND TARGET_QT_RESOURCE_FILES "${FILE}")
      endif()

      # Generated resources are not relative to the target source directory.
      if(NOT RESOURCE_GENERATED)
        cmake_path(RELATIVE_PATH FILE BASE_DIRECTORY ${RELATIVE_SOURCE_DIR} OUTPUT_VARIABLE RESOURCE_FILE_NAME)
      else()
        unset(RESOURCE_GENERATED)
      endif()

      # Add prefix to resource.
      if(HAS_PREFIX)
        set(RESOURCE_ALIAS "/${PREFIX}/${RESOURCE_FILE_NAME}")
      else()
        set(RESOURCE_ALIAS "/${RESOURCE_FILE_NAME}")
      endif()

      # print_message(STATUS "${RESOURCE_ALIAS}")

      set_source_files_properties(${FILE} PROPERTIES "QT_RESOURCE_ALIAS" "${RESOURCE_ALIAS}")
    endforeach()
  endforeach()
endmacro()

# Initialize target
# -------------------------------------------------------------------------------------------------------
macro(snake_init_target TARGET TARGET_PATH LINK_TYPE TARGET_TYPE TARGET_DESCRIPTION TARGET_EXPORTED)
  unset(TARGET_QML_FILES)
  unset(TARGET_QT_TRANSLATION_FILES)
  unset(TARGET_QT_RESOURCE_FILES)

  print_value("  " "PATH" "${TARGET_PATH}")
  print_value("  " "EXPORT" "${TARGET_EXPORTED}")

  set(TARGET_SOURCE_DIR "${CMAKE_SOURCE_DIR}/${TARGET_PATH}")

  if(NOT EXISTS "${TARGET_SOURCE_DIR}")
    message(FATAL_ERROR "Directory does not exist: ${TARGET_PATH}")
  endif()

  file(GLOB_RECURSE SOURCES "${TARGET_SOURCE_DIR}/*.cc")
  file(GLOB_RECURSE HEADERS "${TARGET_SOURCE_DIR}/*.h")

  # Add include directories...
  target_include_directories(${TARGET}
                             ${LINK_TYPE}
                             $<INSTALL_INTERFACE:include>
                             $<BUILD_INTERFACE:${CMAKE_SOURCE_DIR}/lib>
                             $<BUILD_INTERFACE:${CMAKE_BINARY_DIR}/include>
  )

  if(${TARGET_TYPE} STREQUAL "header-library")
    # Need to set the linker language to avoid ambiguity.
    set_target_properties(${TARGET} PROPERTIES LINKER_LANGUAGE CXX)
  elseif(${TARGET_TYPE} STREQUAL "static-library" OR ${TARGET_TYPE} STREQUAL "shared-library")
    if(NOT BUILD_SHARED_LIBS)
      target_compile_definitions(${TARGET} PUBLIC ${TARGET}_STATIC)
    endif()

    # Need to set linker language when directory is empty.
    set_target_properties(${TARGET} PROPERTIES LINKER_LANGUAGE CXX)

    # Install headers and library targets.
    if(${TARGET_EXPORTED})
      cmake_path(RELATIVE_PATH
                 TARGET_SOURCE_DIR
                 BASE_DIRECTORY
                 "${CMAKE_SOURCE_DIR}/lib"
                 OUTPUT_VARIABLE
                 RELATIVE_SOURCE_DIR
      )

      set(EXPORT_HEADER_DIR "${SNAKE_GENERATED_INCLUDE_DIR}/${RELATIVE_SOURCE_DIR}")
      set(EXPORT_HEADER "${EXPORT_HEADER_DIR}/exports.h")

      include(GenerateExportHeader)

      # WARNING: This is internal CMake API; could break in the future. Overrides
      # the internal template used for header generation.
      set(_GENERATE_EXPORT_HEADER_MODULE_DIR "${SNAKE_TEMPLATE_DIR}")

      string(TOUPPER ${SNAKE_PROJECT_NAME} EXPORT_PREFIX)
      string(TOUPPER ${TARGET} EXPORT_NAME)

      if(SNAKE_ENABLE_EXPORT_PREFIX)
        generate_export_header(${TARGET}
                               BASE_NAME
                               "${EXPORT_NAME}"
                               PREFIX_NAME
                               "${EXPORT_PREFIX}_"
                               EXPORT_MACRO_NAME
                               "${EXPORT_NAME}_EXPORT"
                               EXPORT_FILE_NAME
                               "${EXPORT_HEADER}"
        )
      else()
        generate_export_header(${TARGET} BASE_NAME "${EXPORT_NAME}" EXPORT_FILE_NAME "${EXPORT_HEADER}")
      endif()

      list(APPEND HEADERS "${EXPORT_HEADER}")

      # install(TARGETS ${TARGET}
      #         EXPORT ${SNAKE_PROJECT_NAME}-targets
      #         LIBRARY DESTINATION ${SNAKE_PROJECT_LIBRARY_DIRECTORY}
      #         ARCHIVE DESTINATION ${SNAKE_PROJECT_ARCHIVE_DIRECTORY}
      # )

      # cmake_path(RELATIVE_PATH TARGET_PATH BASE_DIRECTORY "${CMAKE_SOURCE_DIR}/lib" OUTPUT_VARIABLE INCLUDE_REL_PATH)

      # install(FILES ${HEADERS} DESTINATION "include/${INCLUDE_REL_PATH}")
    endif()
  elseif(${TARGET_TYPE} STREQUAL "plugin")
    # Faster to find a plugin when it has neither prefixes nor suffixes
    set_target_properties(${TARGET} PROPERTIES SUFFIX "" PREFIX "")

    set_target_properties(${TARGET}
                          PROPERTIES LIBRARY_OUTPUT_DIRECTORY "${CMAKE_MODULES_OUTPUT_DIRECTORY}"
                                     LIBRARY_ARCHIVE_DIRECTORY "${CMAKE_MODULES_OUTPUT_DIRECTORY}"
    )

    target_compile_definitions(${TARGET} PRIVATE QT_PLUGIN QT_DEPRECATED_WARNINGS)

    if(NOT BUILD_SHARED_LIBS)
      print_message(STATUS "Building static plugin...")
      target_compile_definitions(${TARGET} PRIVATE QT_STATICPLUGIN)
    endif()
    
    target_include_directories(${TARGET} PRIVATE $<BUILD_INTERFACE:${TARGET_SOURCE_DIR}>)
  elseif(${TARGET_TYPE} STREQUAL "executable" OR ${TARGET_TYPE} STREQUAL "application")
    cmake_path(RELATIVE_PATH TARGET_SOURCE_DIR BASE_DIRECTORY "${CMAKE_SOURCE_DIR}" OUTPUT_VARIABLE RELATIVE_SOURCE_DIR)

    # Configure config.h file for executables
    set(TARGET_NAME "${TARGET}")
    set(TARGET_DESCRIPTION "${TARGET_DESCRIPTION}")

    configure_file("${SNAKE_TEMPLATE_DIR}/version.h.in"
                   "${SNAKE_GENERATED_INCLUDE_DIR}/${RELATIVE_SOURCE_DIR}/version.h"
                   @ONLY
    )

    target_include_directories(${TARGET}
                               PRIVATE "${SNAKE_GENERATED_INCLUDE_DIR}/${RELATIVE_SOURCE_DIR}"
                                       $<BUILD_INTERFACE:${TARGET_SOURCE_DIR}>
    )
  endif()

  if(SNAKE_ENABLE_QML_MODULES)
    foreach(HEADER ${HEADERS})
      cmake_path(GET HEADER PARENT_PATH HEADER_DIR)
      target_include_directories(${TARGET} ${LINK_TYPE} ${HEADER_DIR})
    endforeach()
  endif()

  list(APPEND ALL_CXX_SOURCES ${HEADERS} ${SOURCES})
  target_sources(${TARGET} ${LINK_TYPE} ${SOURCES} ${HEADERS})
endmacro()

# Helper to find packages (pkg-config).
# -------------------------------------------------------------------------------------------------------
macro(snake_pkg_config NAME)
  find_package(PkgConfig)

  pkg_check_modules(PKG_${NAME} QUIET ${NAME})

  find_path(${NAME}_INCLUDE_DIRS NAMES ${NAME}.hxx ${NAME}.h PATH_SUFFIXES ${NAME} HINTS ${PKG_${NAME}_INCLUDE_DIRS})
  find_library(${NAME}_LIBRARIES NAMES ${PKG_${NAME}_LIBRARIES} ${NAME} HINTS ${PKG_${NAME}_LIBRARY_DIRS})

  include(FindPackageHandleStandardArgs)

  find_package_handle_standard_args(${NAME}
                                    REQUIRED_VARS ${NAME}_LIBRARIES ${NAME}_INCLUDE_DIRS
                                    VERSION_VAR PKG_${NAME}_VERSION
  )
endmacro()

# Finalize target
# -------------------------------------------------------------------------------------------------------
macro(snake_fini_target TARGET)
  get_property(TARGET_LINK_TYPE TARGET ${TARGET} PROPERTY "LINK_LIBRARIES")
  print_value("  " "DEPENDENCIES" "${TARGET_LINK_TYPE}")

  get_property(TARGET_CMAKE_TYPE TARGET ${TARGET} PROPERTY "TYPE")
  print_value("  " "TYPE" "${TARGET_CMAKE_TYPE}")

  # get_property(TARGET_SOURCES TARGET ${TARGET} PROPERTY "SOURCES")
  # list(LENGTH TARGET_SOURCES TARGET_SOURCES_COUNT)
  # print_value("  " "SOURCES" "${TARGET_SOURCES_COUNT}")

  if(TARGET_QT_TRANSLATION_FILES AND SNAKE_ENABLE_TRANSLATIONS)
    find_package(Qt6 REQUIRED COMPONENTS LinguistTools)
    qt_add_translations(${TARGET} ${QT_TRANSLATION_FILES} QM_FILES_OUTPUT_VARIABLE COMPILED_TRANSLATION_FILES)
    install(FILES ${COMPILED_TRANSLATION_FILES} DESTINATION "${CMAKE_INSTALL_PREFIX}/${SNAKE_TRANSLATIONS_DIR}")
  endif()

  if(TARGET_QT_RESOURCE_FILES)
    if(SNAKE_ENABLE_DEVELOPMENT_MODE)
      set(COMPRESSION_LEVEL 1)
    else()
      set(COMPRESSION_LEVEL 19)
    endif()

    qt_add_resources(${TARGET}
                     "${TARGET}"
                     PREFIX
                     "/"
                     FILES
                     ${TARGET_QT_RESOURCE_FILES}
                     OPTIONS
                     -threshold
                     75
                     -compress
                     ${COMPRESSION_LEVEL}
    )
  endif()

  if(TARGET_QML_FILES AND SNAKE_ENABLE_QML_MODULES)
    qt_target_qml_sources(${TARGET} QML_FILES ${TARGET_QML_FILES})
  endif()

  install(TARGETS ${TARGET}
          LIBRARY DESTINATION "${SNAKE_LIBRARY_DIRECTORY}"
          ARCHIVE DESTINATION "${SNAKE_ARCHIVE_DIRECTORY}"
          BUNDLE DESTINATION "${SNAKE_RUNTIME_DIRECTORY}"
          RUNTIME DESTINATION "${SNAKE_RUNTIME_DIRECTORY}"
  )
endmacro()

# We need to wait until QML modules _really_ improves before supporting it. Subscribe to these:
#
# - QTBUG-99653
# - QTBUG-99768
# -------------------------------------------------------------------------------------------------------
macro(snake_add_qml_module TARGET MODULE_URI MODULE_VERSION MODULE_PREFIX)
  if(SNAKE_ENABLE_QML_MODULES)
    get_property(TYPE_PROPERTY TARGET ${TARGET} PROPERTY "TYPE")

    if(TYPE_PROPERTY STREQUAL "EXECUTABLE")
      set(QML_FLAGS NO_PLUGIN NO_RESOURCE_TARGET_PATH)
    else()
      set(QML_FLAGS NO_PLUGIN NO_CREATE_PLUGIN_TARGET NO_GENERATE_PLUGIN_SOURCE)
    endif()

    qt_add_qml_module(${TARGET}
                      ${QML_FLAGS}
                      URI
                      ${MODULE_URI}
                      VERSION
                      ${MODULE_VERSION}
                      RESOURCE_PREFIX
                      "/"
                      IMPORT_PATH
                      ${QT_QML_OUTPUT_DIRECTORY}
    )
  endif()
endmacro()

# Create a graphical application
# -------------------------------------------------------------------------------------------------------
macro(snake_create_graphical_app TARGET)
  if(WINDOWS)
    qt_add_executable(${TARGET} WIN32)
  elseif(OSX)
    qt_add_executable(${TARGET} MACOSX_BUNDLE)
  else()
    qt_add_executable(${TARGET})
  endif()

  set_target_properties(${TARGET}
                        PROPERTIES MACOSX_BUNDLE_GUI_IDENTIFIER ${TARGET}
                                   MACOSX_BUNDLE_BUNDLE_VERSION ${CMAKE_PROJECT_VERSION}
                                   MACOSX_BUNDLE_SHORT_VERSION_STRING
                                   ${CMAKE_PROJECT_VERSION_MAJOR}.${CMAKE_PROJECT_VERSION_MINOR}
  )
endmacro()

# Import a plugin
# -------------------------------------------------------------------------------------------------------
macro(snake_import_plugin TARGET PLUGIN)
  if(BUILD_SHARED_LIBS)
    add_dependencies(${TARGET} ${PLUGIN})
  else()
    target_link_libraries(${TARGET} PRIVATE ${PLUGIN})
  endif()
endmacro()

# Generate a pkg-config module.
# -------------------------------------------------------------------------------------------------------
macro(snake_fetch_pkg NAME)
  print_message(STATUS "Fetching... ${NAME}")
  if(SNAKE_FORCE_UPDATE)
    file(WRITE "${SNAKE_DIR}/.cmake/Find${NAME}.cmake" "snake_pkg_config(${NAME})")
  endif()
  set(${NAME}_FOUND on)
endmacro()

# Fetch git repository.
# -------------------------------------------------------------------------------------------------------
macro(snake_fetch_git NAME URL TAG)
  print_message(STATUS "Fetching... ${NAME}")
  include(FetchContent)
  fetchcontent_declare(${NAME} GIT_REPOSITORY ${URL} GIT_TAG ${TAG})
  fetchcontent_makeavailable(${NAME})
  set(${NAME}_FOUND on)
endmacro()

# Fetch url resource.
# -------------------------------------------------------------------------------------------------------
macro(snake_fetch_url NAME URL HASH)
  print_message(STATUS "Fetching... ${NAME}")
  include(FetchContent)
  fetchcontent_declare(${NAME} GIT_REPOSITORY ${URL} HASH ${HASH})
  fetchcontent_makeavailable(${NAME})
  set(${NAME}_FOUND on)
endmacro()
