cmake_minimum_required(VERSION 3.20)
project(hackAssembler)

set(CMAKE_CXX_STANDARD 20)

include_directories(include)

add_executable(assembler src/main.cpp src/Assembler.cpp include/Assembler.hpp)
