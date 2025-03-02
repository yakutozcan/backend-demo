CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          recipient VARCHAR(15) NOT NULL,
                          content VARCHAR(250) NOT NULL,
                          sent BOOLEAN DEFAULT FALSE,
                          created_at TIMESTAMP DEFAULT NOW()
);


INSERT INTO messages (recipient, content, sent, created_at)
VALUES
    ('+905301234567', 'Hello, how are you?', FALSE, '2025-03-02 08:30:00'),
    ('+905402345678', 'Reminder: Meeting at 3 PM today.', TRUE, '2025-03-02 09:00:00'),
    ('+905503456789', 'Your order has been shipped!', FALSE, '2025-03-02 09:15:00'),
    ('+905604567890', 'Dont forget to submit your report.', TRUE, '2025-03-02 09:30:00'),
  ('+905705678901', 'Happy Birthday! Hope you have a great day.', FALSE, '2025-03-02 10:00:00'),
  ('+905806789012', 'Your subscription has been renewed successfully.', TRUE, '2025-03-02 10:30:00'),
  ('+905907890123', 'Please confirm your attendance for the event.', FALSE, '2025-03-02 11:00:00'),
  ('+905908901234', 'Welcome to the team!', FALSE, '2025-03-02 11:15:00'),
  ('+905909012345', 'Can you send me the report by the end of the day?', TRUE, '2025-03-02 12:00:00'),
  ('+905910123456', 'Your payment is due tomorrow.', FALSE, '2025-03-02 12:30:00'),
  ('+905911234567', 'The server is down. Please check.', FALSE, '2025-03-02 13:00:00'),
  ('+905912345678', 'Meeting rescheduled to 4 PM.', TRUE, '2025-03-02 13:30:00'),
  ('+905913456789', 'Your application has been reviewed.', FALSE, '2025-03-02 14:00:00'),
  ('+905914567890', 'Reminder: Submit your timesheet by 5 PM.', TRUE, '2025-03-02 14:30:00'),
  ('+905915678901', 'Thank you for your feedback!', FALSE, '2025-03-02 15:00:00'),
  ('+905916789012', 'The link to the document is now available.', TRUE, '2025-03-02 15:30:00'),
  ('+905917890123', 'Please reach out if you need any assistance.', FALSE, '2025-03-02 16:00:00'),
  ('+905918901234', 'Your new password is ready to be set.', FALSE, '2025-03-02 16:30:00'),
  ('+905919012345', 'We appreciate your support!', TRUE, '2025-03-02 17:00:00'),
  ('+905920123456', 'Your account has been successfully created.', FALSE, '2025-03-02 17:30:00');
