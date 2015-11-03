def parse_input(raw_input)
    integers = [];
    raw_input.each{ |current_raw|
        begin
            integers.push(Integer(current_raw, 10));
        rescue ArgumentError => detail
            STDERR.puts "An integer could not be parsed. Discarding and continuing. Integer: "+current_raw
        end
    }
    return integers;
end

def main(raw_input)
    #puts "Input:\t"+raw_input.join(", ");

    integers = parse_input(raw_input);
    output = {};
    # structure of output
    # {
    #   integer => {factor => true, factor => true, ...],
    #   ...
    # }
    # the 'true' value is meaningless, it is just a hashmap for free dedupe

    integers.each{ |current_int|
        #puts "Current int:\t"+current_int.to_s;
        output[current_int] = {};
        
        output.each{ |previous_int, factors|
            if(current_int == previous_int)
                break;
            end

            #puts "Previous int:\t"+previous_int.to_s
            if(current_int % previous_int == 0)
                #puts "["+current_int.to_s+"]["+previous_int.to_s+"] = true";
                output[current_int][previous_int] = true;
            end
            if(previous_int % current_int == 0)
                #puts "["+previous_int.to_s+"]["+current_int.to_s+"] = true";
                output[previous_int][current_int] = true;
            end
        }
    }

    return output
end

puts main(ARGV).inspect;
