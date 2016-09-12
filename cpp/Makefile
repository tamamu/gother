
CXX = g++
CXXFLAGS = -g -Wall -Winit-self -Wno-missing-field-initializers
LIBS = -lstdc++fs
TARGET = ./bin/encode
SRCS = $(shell ls *.cpp)
OBJS = $(SRCS:.cpp=.o)
HEADS = $(shell ls *.hpp)

$(TARGET): $(OBJS) $(HEADS)
	$(CXX) -o $(TARGET) $(OBJS) $(LIBS)

all: $(TARGET)

depend:
	$(CXX) -MM -MG $(SRCS) > Makefile.depend
	cat Makefile.depend

clean:
	rm -f $(OBJS) $(TARGET) *~\#*\#

-include Makefile.depend
