PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE MUSCLE(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT NOT NULL DEFAULT '', --肌肉名称
        zh_name TEXT NOT NULL DEFAULT '', --肌肉名称
        overview TEXT NOT NULL DEFAULT '', --简要说明
        tags TEXT NOT NULL DEFAULT '', --标签;逗号分割
        features TEXT NOT NULL DEFAULT '{}' --功能
    
    );
INSERT INTO MUSCLE VALUES(1,'Brachialis','肱肌','起点为肱骨。止点为尺骨冠突。','手臂,肘屈肌','[{"title":"独立功能-向心收缩 ","details":"肘关节屈曲"},{"title":"整体功能-离心收缩","details":"肘关节伸展"},{"title":"整体功能-等长收缩","details":"稳定肘关节"}]');
INSERT INTO MUSCLE VALUES(2,'Biceps brachii','肱二头肌','起点有两个，短头：喙突；长头：肩胛骨盂上结节。止点为桡骨粗隆。','手臂,肘屈肌','[{"title":"独立功能-向心收缩 ","details":"肘关节屈曲；桡尺关节旋后；肩关节屈曲"},{"title":"整体功能-离心收缩","details":"肘关节伸展；桡尺关节旋前；肩关节伸展"},{"title":"整体功能-等长收缩","details":"稳定肘关节和肩胛带"}]');
INSERT INTO MUSCLE VALUES(3,'Triceps brachii','肱三头肌','起点有三个，长头：肩胛骨的孟下结节；外侧头：桡神经沟外上方；内侧头：桡神经沟内下方。止点为尺骨鹰嘴。','手臂,肘伸肌','[{"title":"独立功能-向心收缩","details":"肘关节伸展，肩关节伸展"},{"title":"整体功能-离心收缩","details":"肘关节屈曲，肩关节屈曲"},{"title":"整体功能-等长收缩","details":"稳定肘关节和肩胛带"}]');
INSERT INTO MUSCLE VALUES(4,'Deltoid','三角肌','肩部肌肉群，由三角肌前束、三角肌中束和三角肌后束组成。','肩膀','[{"title":"主要功能","details":"三部分功能各不相同，前束可以让肩关节屈和内旋。中束可以让肩关节外展。后束可以让肩关节伸和外旋。"}]');
INSERT INTO MUSCLE VALUES(5,'Trapezius','斜方肌','连接颈部与背部，由斜方肌上束、斜方肌中束和斜方肌下束组成。','背部','[{"title":"主要功能","details":"斜方肌三部分功能各不相同，上束可以让颈椎伸展、侧屈、旋转，以及肩胛骨上提。中束可以让肩胛后缩。下束可以让肩胛骨下降。"}]');
INSERT INTO MUSCLE VALUES(6,'Pectoralis major','胸大肌','起点为锁骨的前表面，胸骨的前表面，第 1~6 肋软骨。止点为肱骨的大结节。','胸部','[{"title":"独立功能-向心收缩","details":"肩关节屈曲（锁骨纤维）、水平内收和内旋"},{"title":"整体功能-离心收缩","details":"肩关节伸展、水平外展和外旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(7,'Latissimus dorsi','背阔肌','起点为T7~T12 的棘突，骨盆的髂棘，胸腰筋膜，第 9~12 肋骨。止点为肩胛下角；肱骨结节间沟。','背部','[{"title":"独立功能-向心收缩","details":"肩关节伸展、内收和内旋"},{"title":"整体功能-离心收缩","details":"肩关节屈曲、外展和外旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体以及肩部"}]');
INSERT INTO MUSCLE VALUES(8,'Gluteus maximus','臀大肌','起点为髂骨外侧，骶骨和尾骨后侧，骶结节韧带和骶髂后韧带的一部分。止点为臀肌粗隆及髂胫束。','臀部','[{"title":"独立功能-向心收缩","details":"加速髋关节伸展和外旋"},{"title":"整体功能-离心收缩","details":"减缓髋关节屈曲和内旋，并通过髂胫束减缓胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(9,'Quadriceps','股四头肌','腿部的主要肌肉群，由股直肌、股中肌、股外侧肌、股内侧肌四块肌肉组成。','腿部','[{"title":"主要功能","details":"四块肌肉各有其功能，但以膝伸为主。"}]');
INSERT INTO MUSCLE VALUES(12,'Abdominals','腹直肌','起点为耻骨联合。止点为第 5~7 肋骨。','腹部','[{"title":"独立功能-向心收缩","details":"脊柱屈曲、侧屈"},{"title":"整体功能-离心收缩","details":"脊柱伸展、侧屈"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(13,'Obliquus externus abdominis','腹外斜肌','起点为第 5~12 肋骨外表面。止点为髂嵴前部，白线，对侧腹直肌鞘。','腹部','[{"title":"独立功能-向心收缩","details":"脊柱屈曲、侧屈和对侧旋转"},{"title":"整体功能-离心收缩","details":"脊柱伸展、侧屈和旋转"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(14,'Serratus anterior','前锯肌','起点为第 1-9 肋骨。止点为肩胛骨内侧缘。','背部','[{"title":"独立功能-向心收缩","details":"肩胛骨前伸"},{"title":"整体功能-离心收缩","details":"肩胛骨后缩"},{"title":"整体功能-等长收缩","details":"稳定肩胛骨"}]');
INSERT INTO MUSCLE VALUES(15,'Rhomboids','菱形肌','起点为 C6~T4 的棘突。止点为肩胛骨内侧缘。','背部','[{"title":"独立功能-向心收缩","details":"肩胛骨后缩和下回旋"},{"title":"整体功能-离心收缩","details":"肩胛骨前伸和上回旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛骨"}]');
INSERT INTO MUSCLE VALUES(16,'Erector spinae','竖脊肌','竖脊肌由髂肋肌、最长肌、棘肌组成。三块肌肉有共同的起点，骨盆的髂脊、骶骨，止点各不相同。','背部','[{"title":"主要功能","details":"由于三块肌肉基本上是围绕脊柱，所以功能以脊柱的伸为主。"}]');
INSERT INTO MUSCLE VALUES(17,'Sartorius','梨状肌','起点为骶骨前面。止点为股骨大转子。','臀部','[{"title":"独立功能-向心收缩","details":"加速髋关节外旋、外展和伸展"},{"title":"整体功能-离心收缩","details":"减缓髋关节内旋、内收和屈曲"},{"title":"整体功能-等长收缩","details":"稳定髋部和骶髂关节"}]');
INSERT INTO MUSCLE VALUES(18,'Brachioradialis','肱桡肌','起点为肱骨外上髁。止点为桡骨茎突。','手臂,肘屈肌','[{"title":"独立功能-向心收缩 ","details":"肘关节屈曲"},{"title":"整体功能-离心收缩","details":"肘关节伸展"},{"title":"整体功能-等长收缩","details":"稳定肘关节"}]');
INSERT INTO MUSCLE VALUES(19,'Tibialis anterior','胫骨前肌',replace('起点\n·外侧髁及胫骨近端侧面的三分之二\n止点\n·内侧楔⾻内侧面和第一跖骨底','\n',char(10)),'小腿','[{"title":"独立功能-向心收缩","details":"踝关节背屈和足内翻"},{"title":"整体功能-离心收缩","details":"踝关节跖屈和足外翻"},{"title":"整体功能-等长收缩","details":"稳定足弓"}]');
INSERT INTO MUSCLE VALUES(20,'Tibialis posterior','胫骨后肌','起点为腓骨和胫骨近端后面的三分之二。止点为舟骨粗隆，楔骨。','小腿','[{"title":"独立功能-向心收缩","details":"踝关节跖屈和足内翻"},{"title":"整体功能-离心收缩","details":"踝关节背屈和外翻"},{"title":"整体功能-等长收缩","details":"稳定足弓"}]');
INSERT INTO MUSCLE VALUES(21,'Soleus','比目鱼肌','起点为腓骨和胫骨的后表面上部。止点为跟骨跟腱。','小腿','[{"title":"独立功能-向心收缩","details":"加速踝关节跖屈"},{"title":"整体功能-离心收缩","details":"减缓踝关节背屈"},{"title":"整体功能-等长收缩","details":"稳定足部和踝部肌群"}]');
INSERT INTO MUSCLE VALUES(22,'Gastrocnemius','腓肠肌','起点为胫骨、腓骨上端的后面。止点为跟骨跟腱。功能同比目鱼肌。','小腿','[{"title":"独立功能-向心收缩","details":"加速踝关节跖屈"},{"title":"整体功能-离心收缩","details":"减缓踝关节背屈"},{"title":"整体功能-等长收缩","details":"稳定足部和踝部肌群"}]');
INSERT INTO MUSCLE VALUES(23,'Peroneus longus','腓骨长肌','起点为腓骨近端外侧表面的三分之二。止点为楔骨内侧和第一跖骨侧面。','小腿','[{"title":"独立功能-向心收缩","details":"足跖屈、外翻"},{"title":"整体功能-离心收缩","details":"减缓踝关节背屈且使足内翻"},{"title":"整体功能-等长收缩","details":"稳定足部和踝部肌群"}]');
INSERT INTO MUSCLE VALUES(24,'biceps femoris long head','股二头肌长头','起点为坐骨结节，部分骶结节韧带。止点为腓骨头。','大腿','[{"title":"独立功能-向心收缩","details":"加速屈膝、伸髋和胫骨外旋"},{"title":"整体功能-离心收缩","details":"减缓伸膝、屈髋和胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(25,'biceps femoris short head','股二头肌短头','起点为股骨后部的上三分之一处。止点为腓骨头。','大腿','[{"title":"独立功能-向心收缩","details":"加速屈膝、伸髋和胫骨外旋"},{"title":"整体功能-离心收缩","details":"减缓伸膝、屈髋和胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定膝关节"}]');
INSERT INTO MUSCLE VALUES(26,'Semimembranosus','半膜肌','起点为坐骨结节，止点为胫骨内侧髁后面。','大腿','[{"title":"独立功能-向心收缩","details":"加速屈膝、伸髋和胫骨外旋"},{"title":"整体功能-离心收缩","details":"减缓伸膝、屈髋和胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(27,'Semitendinosus','半腱肌','起点为骨盆的坐骨结节和骶结节韧带的一部分。止点为胫骨内侧髁上面（鹅足肌腱位置）。','大腿','[{"title":"独立功能-向心收缩","details":"加速屈膝、伸髋和胫骨外旋"},{"title":"整体功能-离心收缩","details":"减缓伸膝、屈髋和胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(28,'Vastus lateralis','股外侧肌','起点为股骨大转子前侧和下方边缘，臀肌粗隆的外侧面区域，股骨粗线外侧唇。止点为髌骨底部和胫骨粗隆。','大腿','[{"title":"独立功能-向心收缩","details":"加速伸膝"},{"title":"整体功能-离心收缩","details":"减缓屈膝"},{"title":"整体功能-等长收缩","details":"保持膝关节稳定"}]');
INSERT INTO MUSCLE VALUES(29,'Vastus medialis','股内侧肌','起点为股骨粗线内侧唇。止点为髌骨底部，胫骨粗隆。','大腿','[{"title":"独立功能-向心收缩","details":"加速伸膝"},{"title":"整体功能-离心收缩","details":"减缓屈膝"},{"title":"整体功能-等长收缩","details":"保持膝关节稳定"}]');
INSERT INTO MUSCLE VALUES(30,'Vastus intermedius','股中肌','起点为股骨体前外侧上面三分之二区域。止点为髌骨底部，胫骨粗隆。','大腿','[{"title":"独立功能-向心收缩","details":"加速伸膝"},{"title":"整体功能-离心收缩","details":"减缓屈膝"},{"title":"整体功能-等长收缩","details":"保持膝关节稳定"}]');
INSERT INTO MUSCLE VALUES(31,'Rectus femoris','股直肌','起点为髂前下棘，髋臼下缘。止点为髌骨底部，胫骨粗隆。','大腿','[{"title":"独立功能-向心收缩","details":"加速伸膝、屈髋"},{"title":"整体功能-离心收缩","details":"减缓屈膝、伸髋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(32,'Adductor longus','长收肌','起点为耻骨下支前面。止点为股骨粗线中部。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、屈曲"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、伸展"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(33,'Adductor magnus','大收肌（上束）','起点为骨盆的坐骨支。止点为股骨粗线。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、屈曲"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、伸展"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(34,'Adductor magnus','大收肌（下束）','起点为坐骨结节。止点为股骨内上髁。','小腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、伸展"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、屈曲"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(35,'Adductor brevis','短收肌','起点为耻骨下支前面。止点为股骨粗线近端的三分之一。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、屈曲"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、伸展"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(36,'Gracilis','股薄肌','起点为耻骨体下方前面。止点为胫骨近端内侧面（鹅足）。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、屈曲，并辅助胫骨内旋"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、伸展"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(37,'Pectineal','耻骨肌','起点为耻骨上支处。止点为股骨上部后面。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节内收、屈曲"},{"title":"整体功能-离心收缩","details":"减缓髋关节外展、伸展"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(38,'Gluteus Medius','臀中肌（前部）','起点为髂骨翼外面。止点为股骨大转子外面。','臀部','[{"title":"独立功能-向心收缩","details":"加速髋关节外展"},{"title":"整体功能-离心收缩","details":"减缓髋关节内收"},{"title":"整体功能-等长收缩","details":"动态稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(39,'Gluteus Medius','臀中肌（后部）','起点为髂骨翼外面。止点为股骨大转子外面。','臀部','[{"title":"独立功能-向心收缩","details":"加速髋关节外展"},{"title":"整体功能-离心收缩","details":"减缓髋关节内收"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(40,'Gluteus minimus','臀小肌','起点为臀前线与臀下线之间的髂骨处。止点为股骨大转子。','臀部','[{"title":"独立功能-向心收缩","details":"加速髋关节外展"},{"title":"整体功能-离心收缩","details":"减缓髋关节内收"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(41,'Tensor fasciae latae','阔筋膜张肌','起点为髂嵴外面，紧靠髂前上棘后侧。止点为髂胫束近端的三分之一处。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节屈曲、外展和内旋"},{"title":"整体功能-离心收缩","details":"减缓髋关节伸展、内收和外旋等长收缩"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(42,'Psoas major','腰大肌','起点为胸椎最后一块椎骨和腰椎椎骨体（和椎间盘）外侧和横突。止点为股骨小转子。','背部,腰部','[{"title":"独立功能-向心收缩","details":"加速髋关节屈曲和外旋，屈曲和旋转腰椎"},{"title":"整体功能-离心收缩","details":"减缓髋关节伸展和内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(43,'Sartorius','缝匠肌','起点为髂前上棘。止点为胫骨上端内侧面。','大腿','[{"title":"独立功能-向心收缩","details":"加速髋关节屈曲、外旋，加速膝关节屈曲、内旋"},{"title":"整体功能-离心收缩","details":"减缓髋关节伸展、内旋，减缓膝关节伸展、外旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(44,'Obliquus internus abdominis','腹内斜肌','起点为髂脊前三分之二处和胸腰筋膜。止点为第 10~12 肋骨，白线，对侧腹直肌鞘。','腹部','[{"title":"独立功能-向心收缩","details":"脊柱屈曲，侧屈，同侧旋转"},{"title":"整体功能-离心收缩","details":"脊柱伸展、侧屈和旋转"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(45,'Transversus abdominis','腹横肌','起点为第 7~12 肋骨，髂脊的前三分之二处和胸腰筋膜。止点为白线和对侧腹直肌鞘。','腹部','[{"title":"独立功能-向心收缩","details":"增加腹内压，保护腹部内脏"},{"title":"整体功能-等长收缩","details":"与腹内斜肌、多裂肌和深层竖脊肌协同作用来稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(46,'Diaphragm','膈肌',replace('起点有三个部分，\n肋部：第 6~12 肋内侧软骨表面及其相邻的骨质区域。\n胸骨部：剑突的后侧。\n腰部：（1）覆盖了腰方肌和腰大肌外层表面的两个腱膜拱；（2）起于 L1~L3（包含椎间盘）椎体的左右隔脚。\n止点为中心腱。','\n',char(10)),'躯干,腹部','[{"title":"独立功能-向心收缩","details":"膈穹隆下降，胸腔容积扩大"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(47,'Quadratus lumborum','腰方肌','起点为骨盆的髂嵴。止点为第 12 肋骨和 L1~L4 的横突。','背部,腰部','[{"title":"独立功能-向心收缩","details":"脊柱侧屈"},{"title":"整体功能-离心收缩","details":"减缓脊柱向对侧屈曲"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体"}]');
INSERT INTO MUSCLE VALUES(48,'Multifidus','多裂肌','起点为骶骨的后面，腰椎、胸椎和颈椎的骨性突起。止点为起点处上方 1~4 个节段的棘突上。','腰部,背部','[{"title":"独立功能-向心收缩","details":"脊柱伸展和向对侧旋转"},{"title":"整体功能-离心收缩","details":"脊柱的屈曲和旋转"},{"title":"整体功能-等长收缩","details":"稳定脊柱"}]');
INSERT INTO MUSCLE VALUES(49,'Trapezius','斜方肌下束','起点为 T6~T12 的棘突。止点为肩胛冈。','背部','[{"title":"独立功能-向心收缩","details":"肩胛骨下降"},{"title":"整体功能-离心收缩","details":"肩胛骨上提"},{"title":"整体功能-等长收缩","details":"稳定肩胛骨"}]');
INSERT INTO MUSCLE VALUES(50,'Trapezius','斜方肌中束','起点为 T1~T5 棘突。止点为 肩胛骨的肩峰、肩胛冈上方。','背部','[{"title":"独立功能-向心收缩","details":"肩胛骨后缩"},{"title":"整体功能-离心收缩","details":"肩胛骨前伸和上提"},{"title":"整体功能-等长收缩","details":"稳定肩胛骨"}]');
INSERT INTO MUSCLE VALUES(51,'Trapezius','斜方肌上束','起点为枕外隆突、C7 的棘突。止点为锁骨外侧三分之一处、肩胛骨的肩峰。','背部,颈部','[{"title":"独立功能-向心收缩","details":"颈椎伸展、侧屈和旋转；肩胛上提"},{"title":"整体功能-离心收缩","details":"颈椎屈曲、侧屈和旋转；肩胛骨下降"},{"title":"整体功能-等长收缩","details":"稳定颈椎和肩胛骨；在肩胛骨外展和上回旋动作中稳定肩胛骨内侧缘，为原动肌提供稳定的基础"}]');
INSERT INTO MUSCLE VALUES(52,'Musculi patientiae','肩胛提肌','起点为 C1~C4 的横突。止点为肩胛上角。','颈部','[{"title":"独立功能-向心收缩","details":"肩胛骨固定时，颈椎伸展、侧屈和向同侧旋转；颈椎固定时，辅助肩胛骨上提和下回旋"},{"title":"整体功能-离心收缩","details":"肩胛骨固定时，颈椎屈曲、侧屈和对侧旋转；颈椎固定时，肩胛骨下降和上回旋"},{"title":"整体功能-等长收缩","details":"稳定颈椎和肩胛骨"}]');
INSERT INTO MUSCLE VALUES(53,'Pectoralis minor','胸小肌','起点为第 3~5 肋骨。止点为肩胛骨的喙突。','胸部','[{"title":"独立功能-向心收缩","details":"肩胛骨前伸"},{"title":"整体功能-离心收缩","details":"肩胛骨后缩"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(54,'Deltoid','三角肌前束','起点为锁骨的外侧三分之一。止点为肱骨三角肌粗隆。','肩部','[{"title":"独立功能-向心收缩","details":"肩关节屈曲和内旋"},{"title":"整体功能-离心收缩","details":"肩关节伸展和外旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(55,'Deltoid','三角肌中束','起点为肩胛骨的肩峰。止点为肱骨三角肌粗隆。','肩部','[{"title":"独立功能-向心收缩","details":"肩关节外展"},{"title":"整体功能-离心收缩","details":"肩关节内收"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(56,'Deltoid','三角肌后束','起点为肩胛冈。止点为肱骨三角肌粗隆。','肩部','[{"title":"独立功能-向心收缩","details":"肩关节伸展和外旋"},{"title":"整体功能-离心收缩","details":"肩关节屈曲和内旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(57,'Teres minor','小圆肌','起点为肩胛骨的外侧缘。止点为肱骨大结节。','肩袖肌群,背部','[{"title":"独立功能-向心收缩","details":"肩关节外旋"},{"title":"整体功能-离心收缩","details":"肩关节内旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(58,'Infraspinatus','冈下肌','起点为肩胛骨的冈下窝。止点为肱骨大结节的中部。','肩袖肌群,背部','[{"title":"独立功能-向心收缩","details":"肩关节外旋"},{"title":"整体功能-离心收缩","details":"肩关节内旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(59,'Subscapular','肩胛下肌','起点为肩胛骨的肩胛下窝。止点为肱骨小结节。','肩袖肌群,背部','[{"title":"独立功能-向心收缩","details":"肩关节内旋"},{"title":"整体功能-离心收缩","details":"肩关节外旋"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(60,'Supraspinatus','冈上肌','起点为肩胛骨的冈上窝。止点为肱骨大结节的上部。','肩袖肌群,背部','[{"title":"独立功能-向心收缩","details":"肩关节外展"},{"title":"整体功能-离心收缩","details":"肩关节内收"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(61,'Teres major','大圆肌','起点为肩胛骨的肩胛下角。止点为肱骨小结节。','肩袖肌群,背部','[{"title":"独立功能-向心收缩","details":"肩关节内旋、内收和伸展"},{"title":"整体功能-离心收缩","details":"肩关节外旋、外展和屈曲"},{"title":"整体功能-等长收缩","details":"稳定肩胛带"}]');
INSERT INTO MUSCLE VALUES(62,'Anconeus','肘肌','起点为肱骨外上髁。止点为尺骨鹰嘴和尺骨后面。','手臂','[{"title":"独立功能-向心收缩 ","details":"肘关节伸展"},{"title":"整体功能-离心收缩","details":"肘关节屈曲"},{"title":"整体功能-等长收缩","details":"稳定肘关节"}]');
INSERT INTO MUSCLE VALUES(63,'biceps femoris','股二头肌','长头起点为坐骨结节，部分骶结节韧带，止点为腓骨头。短头起点为股骨后部的三分之一处，止点同样为腓骨头。','大腿','[{"title":"独立功能-向心收缩","details":"加速屈膝、伸髋和胫骨外旋"},{"title":"整体功能-离心收缩","details":"减缓伸膝、屈髋和胫骨内旋"},{"title":"整体功能-等长收缩","details":"稳定腰椎-骨盆-髋关节复合体和膝关节"}]');
INSERT INTO MUSCLE VALUES(64,'hamstrings','腘绳肌','大腿后侧主要肌肉群，由股二头肌、半腱肌、半膜肌组成。','大腿','[{"title":"主要功能","details":"三块肌肉各有其功能，但以膝屈为主"}]');
COMMIT;
