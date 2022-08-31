//
// Created by mattl on 8/30/2022.
//

#include "../include/Assembler.hpp"
#include <bitset>

Assembler::Assembler(std::string source, std::string dest) {
    this->sourcePath = source;
    this->destPath = dest;
    this->currIns = 0;

    initPredefined();

    readIntoBuffer();

    dropExtras();

    processLabels();

    convertToBin();

    writeBuffer();

}

std::string Assembler::aToBin(std::string aIns) {
    int num = std::stoi(aIns.substr(1, aIns.size() - 1));

    return std::string("0") + std::bitset<16>(num).to_string().substr(1, 15);
}

std::string Assembler::cToBin(std::string cIns) {
    // ============ comp ============
    // Check if jump is contained
    bool hasJump = (cIns.find(";") != std::string::npos);

    // Check if equals is contained
    bool hasEquals = (cIns.find("=") != std::string::npos);

    int compStart = 0; // first index we want
    int compEnd = 0; // last index we want

    if (hasJump) {
        compEnd = cIns.find(";") - 1;
    }
    else {
        compEnd = cIns.size() - 1;
    }

    if (hasEquals) {
        compStart = cIns.find("=") + 1;
    }
    else {
        compStart = 0;
    }

    std::string comp = cIns.substr(compStart, compEnd - compStart + 1);

    // ============ dest ============
    std::string destTxt = "null";
    if (hasEquals) {
        destTxt = cIns.substr(0, cIns.find("="));
    }

    // ============ jump ============
    std::string jump = "null";
    if (hasJump) {
        jump = cIns.substr(compEnd + 2, cIns.size() - 1);
    }

    std::string cmd = std::string("111") + binFromComp[comp] + binFromDest[destTxt] + binFromJump[jump];
    return cmd;

}

void Assembler::initPredefined() {
    symbolTable["R0"] = currIns++;
    symbolTable["R1"] = currIns++;
    symbolTable["R2"] = currIns++;
    symbolTable["R3"] = currIns++;
    symbolTable["R4"] = currIns++;
    symbolTable["R5"] = currIns++;
    symbolTable["R6"] = currIns++;
    symbolTable["R7"] = currIns++;
    symbolTable["R8"] = currIns++;
    symbolTable["R9"] = currIns++;
    symbolTable["R10"] = currIns++;
    symbolTable["R11"] = currIns++;
    symbolTable["R12"] = currIns++;
    symbolTable["R13"] = currIns++;
    symbolTable["R14"] = currIns++;
    symbolTable["R15"] = currIns++;
    symbolTable["SP"] = 0;
    symbolTable["LCL"] = 1;
    symbolTable["ARG"] = 2;
    symbolTable["THIS"] = 3;
    symbolTable["THAT"] = 4;
    symbolTable["SCREEN"] = 16384;
    symbolTable["KBD"] = 24576;
}

void Assembler::readIntoBuffer() {
    // Open source.txt file
    source.open(sourcePath);

    if (source.is_open()) {
        // Loop through and append each line to the vector
        std::string line;
        while (getline(source, line)) {
            buffer.push_back(line);
        }
    }

    // Close file
    source.close();

}

void Assembler::dropExtras() {
    // Loop through buffer
    tmpBuffer = buffer;
    buffer.clear();

    for (auto iter = tmpBuffer.begin(); iter != tmpBuffer.end(); iter++) {
        // Remove all whitespace and blank lines
        iter->erase(remove(iter->begin(), iter->end(), ' '), iter->end());
        iter->erase(remove(iter->begin(), iter->end(), '\n'), iter->end());

        // Check if it contains a comment
        int ind;
        if ((ind = iter->find("//")) != std::string::npos) {
            // Loop through and find the first occurrence of the slashes and delete those lines
            *iter = iter->substr(0, ind);
        }
        // Remove blank lines
        if (*iter != std::string("")) {
            buffer.push_back(*iter);
        }
    }
}

void Assembler::processLabels() {
    // First pass to load the symbol table
    // Loop through buffer
    int insNum = 0;
    tmpBuffer = buffer;
    buffer.clear();
    for (auto iter = tmpBuffer.begin(); iter != tmpBuffer.end(); iter++) {
        // Check if it is a label
        if ((*iter)[0] == '(') {
            symbolTable[iter->substr(1,iter->size() - 2)] = insNum;
            continue;
        }
        buffer.push_back(*iter);
        insNum++;
    }

    // Now loop through again and fill in labels and variables
    for (auto iter = buffer.begin(); iter != buffer.end(); iter++) {
        // Check if it is an A instruction
        if ((*iter)[0] == '@') {
            // Make sure there are only digits
            bool isNum = true;
            for (int i = 1; i < iter->size(); i++) {
                if (!std::isdigit((*iter)[i])) {
                    isNum = false;
                    break;
                }
            }

            if (isNum) {
                continue;
            }

            // If it is not a number look to process with the symbol table
            std::string symbol = iter->substr(1, iter->size() - 1);
            if (symbolTable.contains(symbol)) {
                *iter = std::string("@") + std::to_string(symbolTable[symbol]);
            }
            else {
                symbolTable[symbol] = currIns++;
                *iter = std::string("@") + std::to_string(symbolTable[symbol]);
            }

        }

    }

}

void Assembler::convertToBin() {
    tmpBuffer = buffer;
    buffer.clear();

    int i = 0;
    for (const std::string& line: tmpBuffer) {
        if (i == 127) {
            i = 100000;
        }
        if (line[0] == '@') {
            buffer.push_back(aToBin(line));
        }
        else {
            buffer.push_back(cToBin(line));
        }
        i++;
    }
}

void Assembler::writeBuffer() {
    dest.open(destPath);

    for (const std::string& line: buffer) {
        dest << line << "\n";
    }

    dest.close();

}