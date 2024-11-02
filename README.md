I have implemented two APIs as specified in the requirements: Upload Receipts and Get Points.


When a receipt is uploaded, it first validates the request body to ensure it can be correctly parsed into the expected JSON format for a Receipt. It also checks that specific string fields, such as dates and dollar amounts, follow the correct formats.
After validation, an ID is assigned to the receipt, which is returned in the response. Next, the points are calculated based on the details of the receipt. The ID, receipt, and calculated points are stored in a map for easy retrieval. Although storing the full receipt data isn't explicitly required, it makes sense to retain it for potential future enhancements.

For the Get Points endpoint, the pre-calculated points are directly retrieved from the in-memory map, providing a quick and efficient response.