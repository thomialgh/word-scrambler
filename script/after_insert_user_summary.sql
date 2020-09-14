CREATE TRIGGER IF NOT EXISTS word_scrambler.after_insert_user_summary AFTER INSERT
ON TABLE user_answers 
FOR EACH ROW
BEGIN  
    declare @correct int
    declare @wrong int
    SELECT 
        @correct = correct_point,
        @wrong = wrong_point,
    FROM word_scrambler.questions
    WHERE id = NEW.question_id;

    declare @point int
    IF(NEW.status = 'wrong') THEN
    SET @point = 0 - @wrong
    ELSE
    SET @point = @correct
    END IF;


    INSERT INTO word_scrambler.user_summary (user_id, score) VALUES (NEW.user_id, @point) ON DUPLICATE KEY UPDATE score=score+@point;     
END