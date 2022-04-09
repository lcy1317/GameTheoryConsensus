
%k = 1:33
%groupNumber =3.*k +1;
clear;
gN= [4:3:100];
fN = 2*floor(gN./3);
badrate = [0.05 0.1 0.16 0.2 0.3 0.4 0.51];
label = []
for pi = 1:length(badrate)
    Poss = []
    p = badrate(pi);
    for i = 1:length(gN)
        pos = 1;
        for j = 0:fN(i)
            pos = pos - nchoosek(gN(i),j)*p^j*(1-p)^(gN(i)-j);
        end
        Poss = [Poss pos];
    end
    plot(gN,Poss);
    label = [label strcat("恶意设备占比：",num2str(p*100),"%")];
    hold on;
end
legend(label);
xlim([0 100]);
ylim([0 0.4]);
xlabel("分组数");
ylabel("攻击成功概率");
hold off;


% Poss = 1 - (groupNumber)