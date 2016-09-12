#include <iostream>
#include <algorithm>
#include <experimental/filesystem>
namespace fs = std::experimental::filesystem::v1;

std::vector<fs::path> getFilePaths(const std::wstring root) {
	std::vector<fs::path> result;
	fs::path p(root); // root
	std::for_each(fs::recursive_directory_iterator(p), fs::recursive_directory_iterator(),
		[&](const fs::path& p) {
			if (fs::is_regular_file(p)) {
				result.push_back(p);
			}
		});
	return result;
}

int main() {
	auto files = getFilePaths(L".");

	for (fs::path p: files)
		std::wcout << p.wstring() << std::endl;

	// name: bytesN_S
	// N for start index
	// S for length of bytes
	
	// 1. Header size
	int32_t bytes0_4;

	// 2. Amount of files
	int32_t bytes4_4 = (int32_t)files.size();
	
	return 0;
}
