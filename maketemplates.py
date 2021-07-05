import random
from PIL import Image, ImageChops
import os
import time

# Make random the same every time it is ran
random.seed(0)

ROUNDS=10

print("Generating templates with %s rounds..." % ROUNDS)
start_time = time.time()

for template_folder in os.listdir("templates"):
    for template in os.listdir("templates/%s" % (template_folder)):
        clean_template = template[:-4]
        template_name = "templates/%s/%s" % (template_folder, template)
        for _ in range(ROUNDS):
            hue = int(random.random() * 255)
            for sat in range(60, 255, 20): 
                color = (hue, 176, sat)
                image = Image.open(template_name)
                image = ImageChops.multiply(image, Image.new("RGBA", (499, 600), color=color))
                temp_name = "assets/%s/%s-(%s-%s-%s).png" % (template_folder, clean_template, color[0], color[1], color[2])
                image.save(temp_name)
        print("Generated templates for %s" % template_name)
print("Generated templates")

end_time = time.time()
print("Time to generate: %ss" % (end_time - start_time))
