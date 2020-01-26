import base64

from locust import HttpLocust, TaskSet, task, between
from random import randint, choice
import uuid

ids = [
	"03fef6ac-1896-4ce8-bd69-b798f85c6e0b", 
	"3395a43e-2d88-40de-b95f-e00e1502085b", 
	"510a0d7e-8e83-4193-b483-e27e09ddc34d",
	"808a2de1-1aaa-4c25-a9b9-6612e8f29a38",
	"819e1fbf-8b7e-4f6d-811f-693534916a8b",
	"837ab141-399e-4c1f-9abc-bace40296bac",
	"a0a4f044-b040-410d-8ead-4de0446aec7e",
	"d3588630-ad8e-49df-bbd7-3167f7efb246",
	"zzz4f044-b040-410d-8ead-4de0446aec7e"
]

class WebTasks(TaskSet):

	@task(2)
	def look_at_catalogue(self):
		self.client.get("/category.html?page=1&size=3")
		self.client.get("/category.html?page=2&size=3")
		self.client.get("/category.html?page=3&size=3")
		self.client.get("/detail.html?id=" + choice(ids))

	@task(1)
	def load(self):
		username = str(uuid.uuid4())
		password = str(uuid.uuid4())
		
		registerBody = {
			"username":	 username, 
			"password":	 password, 
			"email":	 str(uuid.uuid4()) + "@" + str(uuid.uuid4()) + ".org", 
			"firstName": str(uuid.uuid4()), 
			"lastName":  str(uuid.uuid4())
		}
		self.client.post("/register", json=registerBody)

		addressesBody = {
			"street": str(uuid.uuid4()),
			"number": str(uuid.uuid4()).replace("-", ""),
			"country": str(uuid.uuid4()),
			"city": str(uuid.uuid4()),
			"postcode": str(uuid.uuid4()).replace("-", "")
		}
		self.client.post("/addresses", json=addressesBody)

		cardBody = {
			"longNum": "1234432167899876",
			"expires": "11/23",
			"ccv": "123"
		}
		self.client.post("/cards", json=cardBody)

		auth_string = str(base64.b64encode(('{}:{}'.format(username, password)).encode("utf-8")), "utf-8")
		auth_header = 'Basic {}'.format(auth_string)
		item_id = choice(ids)
		self.client.get("/")
		self.client.get("/login", headers={"Authorization":auth_header})
		self.client.get("/category.html")
		self.client.get("/detail.html?id={}".format(item_id))
		self.client.delete("/cart")
		self.client.post("/cart", json={"id": item_id, "quantity": 1})
		self.client.get("/basket.html")
		self.client.post("/orders")


class Web(HttpLocust):
	task_set = WebTasks
	task_set.wait_time = between(0.1, 0.5)
